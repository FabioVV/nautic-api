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

func InsertBoat(acc *models.CreateBoatRequest) error {
	db := storage.GetDB()

	query := "INSERT INTO boats (model, new_used) VALUES ($1, $2)"

	_, err := db.Exec(query, acc.Model, acc.NewUsed)
	if err != nil {
		return err
	}

	return nil
}

func GetBoats(pagenum string, limitPerPage string, model string, price string, id string, active string) ([]models.Boat, int, error) {
	db := storage.GetDB()

	pagenumber, err := strconv.Atoi(pagenum)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve boats (PG1)")
	}
	limit, err := strconv.Atoi(limitPerPage)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve boats (PG2)")
	}

	offset := (pagenumber - 1) * limit

	var boats []models.Boat

	conds := []string{}
	args := []interface{}{}
	paramCount := 1

	if model != "" {
		conds = append(conds, fmt.Sprintf("B.model ILIKE $%d", paramCount))
		args = append(args, "%"+model+"%")
		paramCount++
	}

	if active != "" {
		conds = append(conds, fmt.Sprintf("B.active = $%d", paramCount))
		args = append(args, active)
		paramCount++
	}

	if id != "" {
		conds = append(conds, fmt.Sprintf("B.id = $%d", paramCount))
		args = append(args, id)
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
	SELECT B.id, B.model, B.selling_price, B.cost, B.itens, B.hours, B.year, B.new_used, B.cab_open, B.capacity, B.night_capacity, B.length,
	B.beam, B.draft, B.weight, B.trim, B.fuel_tank_capacity, B.active,
	B.created_at, B.updated_at

	FROM boats AS B

	%s

	ORDER BY B.id, B.model
	LIMIT $%d OFFSET $%d
	`, where, limitArgPos, offsetArgPos)

	rows, err := db.Query(query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return boats, 0, echo.NewHTTPError(http.StatusNotFound, "Boats not found")
		}
		return boats, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve boats")
	}

	queryTotalRecords := fmt.Sprintf(`
	SELECT COUNT(1)
	FROM boats AS B
	%s
	`, where)

	rowsCount := db.QueryRow(queryTotalRecords, args[:len(args)-2]...) // slice to remove the limit and offset args, they are not needed here
	numRecords := 0
	rowsCount.Scan(&numRecords)

	for rows.Next() {
		var curBoat models.Boat
		rows.Scan(&curBoat.Id, &curBoat.Model, &curBoat.PriceSell,
			&curBoat.Cost, &curBoat.Itens, &curBoat.Hours, &curBoat.Year, &curBoat.NewUsed,
			&curBoat.CabOpen, &curBoat.Capacity, &curBoat.NightCapacity, &curBoat.Length, &curBoat.Beam,
			&curBoat.Draft, &curBoat.Weight, &curBoat.Trim, &curBoat.FuelTankCapactiy, &curBoat.Active, &curBoat.CreatedAt, &curBoat.UpdatedAt)
		boats = append(boats, curBoat)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return boats, numRecords, nil
}
