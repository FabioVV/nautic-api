package repositories

import (
	"database/sql"
	"fmt"
	"nautic/cmd/storage"
	"nautic/cmd/utils"
	"nautic/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetEngines(pagenum string, limitPerPage string, model string, active string) ([]models.Engine, int, error) {
	db := storage.GetDB()

	pagenumber, err := strconv.Atoi(pagenum)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve engines (PG1)")
	}
	limit, err := strconv.Atoi(limitPerPage)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve engines (PG2)")
	}

	offset := (pagenumber - 1) * limit

	var engs []models.Engine

	conds := []string{}
	args := []interface{}{}
	paramCount := 1

	if model != "" {
		conds = append(conds, fmt.Sprintf("model ILIKE $%d", paramCount))
		args = append(args, "%"+model+"%")
		paramCount++
	}

	if active != "" {
		conds = append(conds, fmt.Sprintf("active = $%d", paramCount))
		args = append(args, active)
		paramCount++
	}

	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}

	//append pagination range
	args = append(args, limitPerPage, offset)
	limitArgPos := paramCount
	offsetArgPos := paramCount + 1

	query := fmt.Sprintf(`
	SELECT id, model, type, weight, rotation, power, cylinders, selling_price, command, clocks, tempo, fuel_type, active, created_at, updated_at, propulsion
	FROM engines
	%s
	ORDER BY id, model
	LIMIT $%d OFFSET $%d
	`, where, limitArgPos, offsetArgPos)

	rows, err := db.Query(query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return engs, 0, echo.NewHTTPError(http.StatusNotFound, "Users not found")
		}
		return engs, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve accs"+err.Error())
	}

	queryTotalRecords := fmt.Sprintf(`
	SELECT COUNT(1)
	FROM engines AS E
	%s
	`, where)

	rowsCount := db.QueryRow(queryTotalRecords, args[:len(args)-2]...) // slice to remove the limit and offset args, they are not needed here
	numRecords := 0
	rowsCount.Scan(&numRecords)

	for rows.Next() {
		var curAcc models.Engine

		rows.Scan(&curAcc.Id, &curAcc.Model, &curAcc.Type, &curAcc.Weight, &curAcc.Rotation, &curAcc.Power, &curAcc.Cylinders,
			&curAcc.PriceSell, &curAcc.Command, &curAcc.Clocks, &curAcc.Tempo, &curAcc.FuelType, &curAcc.Active, &curAcc.CreatedAt, &curAcc.UpdatedAt, &curAcc.Propulsion)
		engs = append(engs, curAcc)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return engs, numRecords, nil
}

func InsertEngine(acc *models.CreateEngineRequest) error {
	db := storage.GetDB()

	query := `
INSERT INTO engines
    (model, type, weight, rotation, power, cylinders, selling_price, command, clocks, tempo, fuel_type, propulsion)
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`

	var newID int64
	err := db.QueryRow(
		query,
		acc.Model,
		acc.Type,
		acc.Weight,
		acc.Rotation,
		acc.Power,
		acc.Cylinders,
		acc.PriceSell,
		acc.Command,
		acc.Clocks,
		acc.Tempo,
		acc.FuelType,
		acc.Propulsion,
	).Scan(&newID)

	if err != nil {
		if _, ok := utils.CheckForError("unique_type", err); ok {
			return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"errors": echo.Map{"type": "engine already exists"}})
		}
		return err
	}

	return nil
}

func GetEngine(id int) (models.Engine, error) {
	db := storage.GetDB()

	var curAcc models.Engine
	query := `
	SELECT id, model, type, weight, rotation, power, cylinders, selling_price, command, clocks, tempo, fuel_type, active, created_at, updated_at, propulsion
	FROM engines
	WHERE id = $1`

	if err := db.QueryRow(query, id).Scan(&curAcc.Id, &curAcc.Model, &curAcc.Type, &curAcc.Weight, &curAcc.Rotation, &curAcc.Power, &curAcc.Cylinders,
		&curAcc.PriceSell, &curAcc.Command, &curAcc.Clocks, &curAcc.Tempo, &curAcc.FuelType, &curAcc.Active, &curAcc.CreatedAt, &curAcc.UpdatedAt, &curAcc.Propulsion); err != nil {
		if err == sql.ErrNoRows {
			return curAcc, echo.NewHTTPError(http.StatusNotFound, "Engine not found")
		}
		return curAcc, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve engine")
	}

	return curAcc, nil
}

func DeactivateEngine(id int) error {
	db := storage.GetDB()

	_, err := GetEngine(id)
	if err != nil {
		return err
	}

	query := `UPDATE engines SET active = 'N' WHERE id = $1`

	_, err = db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
