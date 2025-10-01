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

func GetBoat(id int) (models.Boat, error) {
	db := storage.GetDB()

	var curBoat models.Boat

	query := `
	SELECT B.id, B.model, B.selling_price, B.cost, B.itens, B.hours, B.year, B.new_used, B.cab_open, B.capacity, B.night_capacity, B.length,
	B.beam, B.draft, B.weight, B.trim, B.fuel_tank_capacity, B.active,
	B.created_at, B.updated_at

	FROM boats AS B

	WHERE B.id = $1

	ORDER BY B.id, B.model
	`

	if err := db.QueryRow(query, id).Scan(&curBoat.Id, &curBoat.Model, &curBoat.PriceSell,
		&curBoat.Cost, &curBoat.Itens, &curBoat.Hours, &curBoat.Year, &curBoat.NewUsed,
		&curBoat.CabOpen, &curBoat.Capacity, &curBoat.NightCapacity, &curBoat.Length, &curBoat.Beam,
		&curBoat.Draft, &curBoat.Weight, &curBoat.Trim, &curBoat.FuelTankCapactiy, &curBoat.Active, &curBoat.CreatedAt, &curBoat.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return curBoat, echo.NewHTTPError(http.StatusNotFound, "Negotiation not found")
		}
		return curBoat, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve boat")
	}

	return curBoat, nil
}

func UpdateBoat(id int, cT *models.BoatRequest) error {
	db := storage.GetDB()

	query := `UPDATE boats SET `
	params := []interface{}{}
	paramCount := 0

	if cT.Model != nil {
		paramCount++
		query += fmt.Sprintf("model = $%d, ", paramCount)
		params = append(params, *&cT.Model)
	}

	if cT.Beam != nil {
		paramCount++
		query += fmt.Sprintf("beam = $%d, ", paramCount)
		params = append(params, *&cT.Beam)
	}

	if cT.Trim != nil {
		paramCount++
		query += fmt.Sprintf("trim = $%d, ", paramCount)
		params = append(params, *&cT.Trim)
	}

	if cT.Capacity != nil {
		paramCount++
		query += fmt.Sprintf("capacity = $%d, ", paramCount)
		params = append(params, *&cT.Capacity)
	}

	if cT.NightCapacity != nil {
		paramCount++
		query += fmt.Sprintf("night_capacity = $%d, ", paramCount)
		params = append(params, *&cT.NightCapacity)
	}

	if cT.Weight != nil {
		paramCount++
		query += fmt.Sprintf("weight = $%d, ", paramCount)
		params = append(params, *&cT.Weight)
	}

	if cT.Length != nil {
		paramCount++
		query += fmt.Sprintf("length = $%d, ", paramCount)
		params = append(params, *&cT.Length)
	}

	if cT.FuelTankCapactiy != nil {
		paramCount++
		query += fmt.Sprintf("fuel_tank_capacity = $%d, ", paramCount)
		params = append(params, *&cT.FuelTankCapactiy)
	}

	if cT.Draft != nil {
		paramCount++
		query += fmt.Sprintf("draft = $%d, ", paramCount)
		params = append(params, *&cT.Draft)
	}

	if cT.Itens != nil {
		paramCount++
		query += fmt.Sprintf("itens = $%d, ", paramCount)
		params = append(params, *&cT.Itens)
	}

	if cT.PriceSell != nil {
		paramCount++
		query += fmt.Sprintf("selling_price = $%d, ", paramCount)
		params = append(params, *&cT.PriceSell)
	}

	if cT.Hours != nil {
		paramCount++
		query += fmt.Sprintf("hours = $%d, ", paramCount)
		params = append(params, *&cT.Hours)
	}

	if cT.Year != nil {
		paramCount++
		query += fmt.Sprintf("year = $%d, ", paramCount)
		params = append(params, *&cT.Year)
	}

	if cT.NewUsed != nil {
		paramCount++
		query += fmt.Sprintf("new_used = $%d, ", paramCount)
		params = append(params, *&cT.NewUsed)
	}

	if cT.CabOpen != nil {
		paramCount++
		query += fmt.Sprintf("cab_open = $%d, ", paramCount)
		params = append(params, *&cT.CabOpen)
	}

	if len(params) == 0 {
		return nil
	}

	//Remove the trailing comma and space from the query
	query = query[:len(query)-2]

	paramCount++
	query += fmt.Sprintf(" WHERE id = $%d", paramCount)
	params = append(params, id)

	_, err := db.Exec(query, params...)
	if err != nil {
		return err
	}

	return nil
}
