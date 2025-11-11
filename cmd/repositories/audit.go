package repositories

import (
	"database/sql"
	"fmt"
	"nautic/cmd/audit"
	"nautic/cmd/storage"
	"nautic/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetLogs(id_user int, pageNumber string, perPage string, action string, description string) ([]models.Log, int, error) {
	db := storage.GetDB()

	pagenumber, err := strconv.Atoi(pageNumber)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve users (PG1)")
	}
	limit, err := strconv.Atoi(perPage)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve users (PG2)")
	}

	offset := (pagenumber - 1) * limit

	var logs []models.Log

	conds := []string{}
	args := []interface{}{}
	paramCount := 1

	if action != "" {
		conds = append(conds, fmt.Sprintf("L.action ILIKE $%d", paramCount))
		args = append(args, "%"+action+"%")
		paramCount++
	}

	if description != "" {
		conds = append(conds, fmt.Sprintf("L.extra_description ILIKE $%d", paramCount))
		args = append(args, "%"+description+"%")
		paramCount++
	}

	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}

	//append pagination range
	args = append(args, perPage, offset)
	limitArgPos := paramCount
	offsetArgPos := paramCount + 1

	query := fmt.Sprintf(`
	SELECT L.id, U.id, U.name, L.url, L.action, L.extra_description, L.created_at
	FROM logs_2025 AS L
	INNER JOIN users AS U ON U.id = L.id_user
	%s
	ORDER BY L.id, U.name
	LIMIT $%d OFFSET $%d
	`, where, limitArgPos, offsetArgPos)

	rows, err := db.Query(query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return logs, 0, echo.NewHTTPError(http.StatusNotFound, "Logs not found")
		}
		return logs, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve logs"+err.Error())
	}

	queryTotalRecords := fmt.Sprintf(`
	SELECT COUNT(1)
	FROM logs_2025 AS L
	%s
	`, where)
	//println(queryTotalRecords)

	rowsCount := db.QueryRow(queryTotalRecords, args[:len(args)-2]...) // slice to remove the limit and offset args, they are not needed here
	numRecords := 0
	rowsCount.Scan(&numRecords)

	for rows.Next() {
		var curLog models.Log
		rows.Scan(&curLog.Id, &curLog.UserId, &curLog.Username, &curLog.Url, &curLog.Action, &curLog.Description, &curLog.CreatedAt)
		logs = append(logs, curLog)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	_ = audit.InsertLog(id_user, "/logs", "GET", fmt.Sprintf("GET request on logs"), query)

	return logs, numRecords, nil
}
