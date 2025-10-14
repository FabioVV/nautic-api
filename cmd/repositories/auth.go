package repositories

import (
	"database/sql"
	"nautic/cmd/storage"
	"nautic/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetPermissions() ([]models.Permission, error) {
	db := storage.GetDB()

	var perms []models.Permission

	query := `
	SELECT P.id, P.module, P.code, P.description, P.url
	FROM permissions AS P

	ORDER BY P.module
	`

	rows, err := db.Query(query)

	if err != nil {
		if err == sql.ErrNoRows {
			return perms, echo.NewHTTPError(http.StatusNotFound, "perms not found")
		}
		return perms, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve perms")
	}

	for rows.Next() {
		var curPerm models.Permission
		rows.Scan(&curPerm.Id, &curPerm.Module, &curPerm.Code, &curPerm.Description, &curPerm.Url)
		perms = append(perms, curPerm)
	}

	return perms, nil
}
