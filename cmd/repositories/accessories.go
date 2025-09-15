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

func UpdateAccessory(id int, accT *models.UpdateAccessoryRequest) error {
	db := storage.GetDB()

	accTg, err := GetAccessory(id)
	if err != nil {
		return err
	}

	if accTg.Active == "N" {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"errors": echo.Map{"accessory": "Accessory must bet active to update it"}})
	}

	query := `UPDATE accessories SET `
	params := []interface{}{}
	paramCount := 0

	if accT.Model != nil {
		paramCount++
		query += fmt.Sprintf("Model = $%d, ", paramCount)
		params = append(params, *&accT.Model)
	}

	if accT.PriceBuy != nil {
		paramCount++
		query += fmt.Sprintf("price_buy = $%d, ", paramCount)
		params = append(params, *&accT.PriceBuy)
	}

	if accT.PriceSell != nil {
		paramCount++
		query += fmt.Sprintf("price_sell = $%d, ", paramCount)
		params = append(params, *&accT.PriceSell)
	}

	if accT.Details != nil {
		paramCount++
		query += fmt.Sprintf("details = $%d, ", paramCount)
		params = append(params, *&accT.Details)
	}

	if accT.IdAccessoryType != nil {
		paramCount++
		query += fmt.Sprintf("id_accessory_type = $%d, ", paramCount)
		params = append(params, *&accT.IdAccessoryType)
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

func UpdateAccessoryType(id int, accT *models.UpdateAccessoryTypeRequest) error {
	db := storage.GetDB()

	accTg, err := GetAccessoryType(id)
	if err != nil {
		return err
	}

	if accTg.Active == "N" {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"errors": echo.Map{"type": "type must bet active to update it"}})

	}

	query := `UPDATE accessory_types SET `
	params := []interface{}{}
	paramCount := 0

	if accT.Type != nil {
		paramCount++
		query += fmt.Sprintf("type = $%d, ", paramCount)
		params = append(params, *accT.Type)
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

func DeactivateAccessoryType(id int) error {
	db := storage.GetDB()

	_, err := GetAccessoryType(id)
	if err != nil {
		return err
	}

	query := `UPDATE accessory_types SET active = 'N' WHERE id = $1`

	_, err = db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil

}

func DeactivateAccessory(id int) error {
	db := storage.GetDB()

	_, err := GetAccessory(id)
	if err != nil {
		return err
	}

	query := `UPDATE accessories SET active = 'N' WHERE id = $1`

	_, err = db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func GetAccessoryType(id int) (models.AccessoryType, error) {
	db := storage.GetDB()

	var accT models.AccessoryType
	query := `SELECT id, type, active, created_at, updated_at FROM accessory_types WHERE id = $1`

	if err := db.QueryRow(query, id).Scan(&accT.Id, &accT.Type, &accT.Active, &accT.CreatedAt, &accT.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return accT, echo.NewHTTPError(http.StatusNotFound, "Accessory type not found")
		}
		return accT, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve accessory type")
	}

	return accT, nil
}

func GetAccessory(id int) (models.Accessory, error) {
	db := storage.GetDB()

	var acc models.Accessory
	query := `SELECT A.id, A.model, A.details, A.price_buy, A.price_sell, A.active, A.created_at, A.updated_at, A.id_accessory_type, AT.type
	FROM accessories AS A
	INNER JOIN accessory_types AS AT ON A.id_accessory_type = AT.id

	WHERE A.id = $1`

	if err := db.QueryRow(query, id).Scan(&acc.Id, &acc.Model, &acc.Details, &acc.PriceBuy, &acc.PriceSell, &acc.Active, &acc.CreatedAt, &acc.UpdatedAt, &acc.IdAccessoryType, &acc.Type); err != nil {
		if err == sql.ErrNoRows {
			return acc, echo.NewHTTPError(http.StatusNotFound, "Accessory not found")
		}
		return acc, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve accessory" + err.Error())
	}

	return acc, nil
}

func InsertAccessoryType(accT *models.CreateAccessoryTypeRequest) error {
	db := storage.GetDB()

	query := "INSERT INTO accessory_types (type) VALUES ($1)"

	_, err := db.Exec(query, accT.Type)
	if err != nil {
		if _, ok := utils.CheckForUserError("unique_type", err); ok {
			return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"errors": echo.Map{"type": "type already exists"}})
		}
		return err
	}

	return nil
}

func InsertAccessory(acc *models.CreateAccessoryRequest) error {
	db := storage.GetDB()

	query := "INSERT INTO accessories (model, details, price_buy, price_sell, id_accessory_type) VALUES ($1, $2, $3, $4, $5)"

	_, err := db.Exec(query, acc.Model, acc.Details, acc.PriceBuy, acc.PriceSell, acc.IdAccessoryType)
	if err != nil {
		if _, ok := utils.CheckForUserError("unique_type", err); ok {
			return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"errors": echo.Map{"type": "accessory already exists"}})
		}
		return err
	}

	return nil
}

func GetAccessoriesTypes(pagenum string, limitPerPage string, _type string, active string) ([]models.AccessoryType, int, error) {
	db := storage.GetDB()

	pagenumber, err := strconv.Atoi(pagenum)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve accs types (PG1)")
	}
	limit, err := strconv.Atoi(limitPerPage)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve accs types (PG2)")
	}

	offset := (pagenumber - 1) * limit

	var accs []models.AccessoryType

	conds := []string{}
	args := []interface{}{}
	paramCount := 1

	if _type != "" {
		conds = append(conds, fmt.Sprintf("A.type ILIKE $%d", paramCount))
		args = append(args, "%"+_type+"%")
		paramCount++
	}

	if active != "" {
		conds = append(conds, fmt.Sprintf("A.active = $%d", paramCount))
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
	SELECT A.id, A.type, A.created_at, A.updated_at, A.active
	FROM accessory_types AS A
	%s
	ORDER BY A.id, A.type
	LIMIT $%d OFFSET $%d
	`, where, limitArgPos, offsetArgPos)

	rows, err := db.Query(query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return accs, 0, echo.NewHTTPError(http.StatusNotFound, "Types not found")
		}
		return accs, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve accs"+err.Error())
	}

	queryTotalRecords := fmt.Sprintf(`
	SELECT COUNT(1)
	FROM accessory_types AS A
	%s
	`, where)
	//println(queryTotalRecords)

	rowsCount := db.QueryRow(queryTotalRecords, args[:len(args)-2]...) // slice to remove the limit and offset args, they are not needed here
	numRecords := 0
	rowsCount.Scan(&numRecords)

	for rows.Next() {
		var curAcc models.AccessoryType
		rows.Scan(&curAcc.Id, &curAcc.Type, &curAcc.CreatedAt, &curAcc.UpdatedAt, &curAcc.Active)
		accs = append(accs, curAcc)
	}

	return accs, numRecords, nil
}

func GetAccessories(pagenum string, limitPerPage string, model string, active string) ([]models.Accessory, int, error) {
	db := storage.GetDB()

	pagenumber, err := strconv.Atoi(pagenum)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve accs (PG1)")
	}
	limit, err := strconv.Atoi(limitPerPage)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve accs (PG2)")
	}

	offset := (pagenumber - 1) * limit

	var accs []models.Accessory

	conds := []string{}
	args := []interface{}{}
	paramCount := 1

	if model != "" {
		conds = append(conds, fmt.Sprintf("A.model ILIKE $%d", paramCount))
		args = append(args, "%"+model+"%")
		paramCount++
	}

	if active != "" {
		conds = append(conds, fmt.Sprintf("A.active = $%d", paramCount))
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
	SELECT A.id, A.model, A.details, A.price_buy, A.price_sell, A.created_at, A.updated_at, A.active, AT.type
	FROM accessories AS A
	INNER JOIN accessory_types AS AT ON A.id_accessory_type = AT.id
	%s
	ORDER BY A.id, A.model
	LIMIT $%d OFFSET $%d
	`, where, limitArgPos, offsetArgPos)

	rows, err := db.Query(query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return accs, 0, echo.NewHTTPError(http.StatusNotFound, "Users not found")
		}
		return accs, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve accs"+err.Error())
	}

	queryTotalRecords := fmt.Sprintf(`
	SELECT COUNT(1)
	FROM accessories AS A
	%s
	`, where)

	rowsCount := db.QueryRow(queryTotalRecords, args[:len(args)-2]...) // slice to remove the limit and offset args, they are not needed here
	numRecords := 0
	rowsCount.Scan(&numRecords)

	for rows.Next() {
		var curAcc models.Accessory
		rows.Scan(&curAcc.Id, &curAcc.Model, &curAcc.Details, &curAcc.PriceBuy, &curAcc.PriceSell, &curAcc.CreatedAt, &curAcc.UpdatedAt, &curAcc.Active, &curAcc.Type)
		accs = append(accs, curAcc)
	}

	return accs, numRecords, nil
}
