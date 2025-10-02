package repositories

import (
	"database/sql"
	"fmt"
	"nautic/cmd/storage"
	"nautic/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetNegotiationsReport(pagenum string, limitPerPage string, name string, boat string) ([]models.Negotiation, int, error) {
	db := storage.GetDB()

	var negs []models.Negotiation

	pagenumber, err := strconv.Atoi(pagenum)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve boats (PG1)")
	}
	limit, err := strconv.Atoi(limitPerPage)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve boats (PG2)")
	}

	offset := (pagenumber - 1) * limit

	conds := []string{}
	args := []interface{}{}
	paramCount := 1

	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}

	args = append(args, limitPerPage, offset)
	limitArgPos := paramCount
	offsetArgPos := paramCount + 1

	query := fmt.Sprintf(`
	SELECT SB.id, 
			SB.id_customer,
	 		SB.id_mean_communication, 
			C.name,
			C.email,
			C.phone,
			MC.name,
			SB.boat_name, 
			SB.estimated_value, 
			SB.max_estimated_value, 
			SB.customer_city, 
			SB.customer_navigation_city, 
			SB.boat_capacity_needed, 
			SB.new_used, 
			SB.cab_open, 
			SB.stage, 
			C.qualified,
			(SB.created_at < now() - interval '24 hours') AS has_passed_24hrs
	FROM so_business AS SB

	INNER JOIN customers AS C ON SB.id_customer = C.id
	INNER JOIN mean_communication AS MC ON SB.id_mean_communication = MC.id

	%s
	ORDER BY SB.id
	LIMIT $%d OFFSET $%d
	`, where, limitArgPos, offsetArgPos)

	rows, err := db.Query(query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return negs, 0, echo.NewHTTPError(http.StatusNotFound, "Negotiations not found")
		}
		return negs, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve negotiations")
	}

	queryTotalRecords := fmt.Sprintf(`
	SELECT COUNT(1)
	FROM so_business AS SB
	INNER JOIN customers AS C ON SB.id_customer = C.id
	INNER JOIN mean_communication AS MC ON SB.id_mean_communication = MC.id
	%s
	`, where)
	//println(queryTotalRecords)

	rowsCount := db.QueryRow(queryTotalRecords, args...)
	numRecords := 0
	rowsCount.Scan(&numRecords)

	for rows.Next() {
		var curNeg models.Negotiation

		if err := rows.Scan(&curNeg.Id, &curNeg.CustomerId, &curNeg.MeanComId,
			&curNeg.Name, &curNeg.Email, &curNeg.Phone, &curNeg.MeamComName,
			&curNeg.BoatName, &curNeg.EstimatedValue, &curNeg.MaxEstimatedValue, &curNeg.City,
			&curNeg.NavigationCity, &curNeg.BoatCapacityNeeded, &curNeg.NewUsed, &curNeg.CabOpen, &curNeg.Stage, &curNeg.Qualified, &curNeg.HasPassed24Hrs); err != nil {
			return nil, 0, fmt.Errorf("scan error: %w", err)
		}

		negs = append(negs, curNeg)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return negs, numRecords, nil
}
