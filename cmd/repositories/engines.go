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

func UpdateEngine(id int, accT *models.CreateEngineRequest) error {
	db := storage.GetDB()

	accTg, err := GetEngine(id)
	if err != nil {
		return err
	}

	if accTg.Active == "N" {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"errors": echo.Map{"accessory": "Engine must bet active to update it"}})
	}

	query := `UPDATE engines SET `
	params := []interface{}{}
	paramCount := 0

	paramCount++
	query += fmt.Sprintf("Model = $%d, ", paramCount)
	params = append(params, *&accT.Model)

	if accT.PriceSell != nil {
		paramCount++
		query += fmt.Sprintf("selling_price = $%d, ", paramCount)
		params = append(params, *&accT.PriceSell)
	}

	if accT.Type != nil {
		paramCount++
		query += fmt.Sprintf("type = $%d, ", paramCount)
		params = append(params, *&accT.Type)
	}

	if accT.Propulsion != nil {
		paramCount++
		query += fmt.Sprintf("propulsion = $%d, ", paramCount)
		params = append(params, *&accT.Propulsion)
	}

	if accT.Weight != nil {
		paramCount++
		query += fmt.Sprintf("weight = $%d, ", paramCount)
		params = append(params, *&accT.Weight)
	}

	if accT.Rotation != nil {
		paramCount++
		query += fmt.Sprintf("rotation = $%d, ", paramCount)
		params = append(params, *&accT.Rotation)
	}

	if accT.Power != nil {
		paramCount++
		query += fmt.Sprintf("power = $%d, ", paramCount)
		params = append(params, *&accT.Power)
	}

	if accT.Cylinders != nil {
		paramCount++
		query += fmt.Sprintf("cylinders = $%d, ", paramCount)
		params = append(params, *&accT.Cylinders)
	}

	if accT.Command != nil {
		paramCount++
		query += fmt.Sprintf("command = $%d, ", paramCount)
		params = append(params, *&accT.Command)
	}

	if accT.Clocks != nil {
		paramCount++
		query += fmt.Sprintf("clocks = $%d, ", paramCount)
		params = append(params, *&accT.Clocks)
	}

	if accT.Tempo != nil {
		paramCount++
		query += fmt.Sprintf("tempo = $%d, ", paramCount)
		params = append(params, *&accT.Tempo)
	}

	if accT.FuelType != nil {
		paramCount++
		query += fmt.Sprintf("fuel_type = $%d, ", paramCount)
		params = append(params, *&accT.FuelType)
	}

	if len(params) == 0 {
		return nil
	}

	//Remove the trailing comma and space from the query
	query = query[:len(query)-2]

	paramCount++
	query += fmt.Sprintf(" WHERE id = $%d", paramCount)
	params = append(params, id)

	_, err = db.Exec(query, params...)
	if err != nil {
		return err
	}

	return nil
}
