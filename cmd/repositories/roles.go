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

func InsertRole(role *models.CreateRole) error {
	db := storage.GetDB()

	query := "SELECT id FROM roles WHERE name = $1"

	var exists int = 0
	db.QueryRow(query, role.Name).Scan(&exists)

	if exists != 0 {
		return echo.NewHTTPError(http.StatusInternalServerError, "Role already exists")
	} else {
		query := "INSERT INTO roles (name) VALUES ($1)"

		_, err := db.Exec(query, role.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetRoles(pagenum string, limitPerPage string, name string, showAdmin string) ([]models.Role, int, error) {
	db := storage.GetDB()

	pagenumber, err := strconv.Atoi(pagenum)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve roles (PG1)")
	}
	limit, err := strconv.Atoi(limitPerPage)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve roles (PG2)")
	}

	offset := (pagenumber - 1) * limit

	var roles []models.Role

	conds := []string{}
	args := []any{}
	paramCount := 1

	if name != "" {
		conds = append(conds, fmt.Sprintf("R.name ILIKE $%d", paramCount))
		args = append(args, "%"+name+"%")
		paramCount++
	}

	if showAdmin == "N" {
		conds = append(conds, fmt.Sprintf("LOWER(R.name) <> LOWER($%d)", paramCount))
		args = append(args, "admin")
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
	SELECT R.id, R.name, R.created_at, R.updated_at
	FROM roles AS R
	%s
	ORDER BY R.id, R.name
	LIMIT $%d OFFSET $%d
	`, where, limitArgPos, offsetArgPos)
	println(query)

	rows, err := db.Query(query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return roles, 0, echo.NewHTTPError(http.StatusNotFound, "Roles not found")
		}
		return roles, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve roles"+err.Error())
	}

	queryTotalRecords := fmt.Sprintf(`
	SELECT COUNT(1)
	FROM roles AS R
	%s
	`, where)

	rowsCount := db.QueryRow(queryTotalRecords, args[:len(args)-2]...) // slice to remove the limit and offset args, they are not needed here
	numRecords := 0
	rowsCount.Scan(&numRecords)

	for rows.Next() {
		var curRole models.Role
		rows.Scan(&curRole.Id, &curRole.Name, &curRole.CreatedAt, &curRole.UpdatedAt)
		roles = append(roles, curRole)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return roles, numRecords, nil
}

func GetRole(id int) (models.Role, error) {
	db := storage.GetDB()

	var role models.Role
	query := `SELECT id, name, created_at, updated_at FROM roles WHERE id = $1`

	if err := db.QueryRow(query, id).Scan(&role.Id, &role.Name, &role.CreatedAt, &role.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return role, echo.NewHTTPError(http.StatusNotFound, "Role not found")
		}
		return role, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve role"+err.Error())
	}

	return role, nil
}

func DeleteRole(id int) error {
	db := storage.GetDB()

	_, err := GetRole(id)
	if err != nil {
		return err
	}

	query := `DELETE FROM roles WHERE id = $1`

	_, err = db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil

}

func UpdateRole(id int, rM *models.CreateRole) error {
	db := storage.GetDB()

	_, err := GetRole(id)
	if err != nil {
		return err
	}

	query := `UPDATE roles SET `
	params := []any{}
	paramCount := 0

	if rM.Name != nil {
		paramCount++
		query += fmt.Sprintf("name = $%d, ", paramCount)
		params = append(params, *rM.Name)
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
