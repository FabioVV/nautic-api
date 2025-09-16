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

func UpdateComMean(id int, mcR *models.UpdateCommunicationMeaneRequest) error {
	db := storage.GetDB()

	mc, err := GetComMean(id)
	if err != nil {
		return err
	}

	if mc.Active == "N" {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"errors": echo.Map{"name": "mean must bet active to update it"}})

	}

	query := `UPDATE mean_communication SET `
	params := []interface{}{}
	paramCount := 0

	if mcR.Name != nil {
		paramCount++
		query += fmt.Sprintf("name = $%d, ", paramCount)
		params = append(params, *mcR.Name)
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

func DeactivateComMean(id int) error {
	db := storage.GetDB()

	_, err := GetComMean(id)
	if err != nil {
		return err
	}

	query := `UPDATE mean_communication SET active = 'N' WHERE id = $1`

	_, err = db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil

}

func GetComMean(id int) (models.CommunicationMean, error) {
	db := storage.GetDB()

	var mc models.CommunicationMean
	query := `SELECT id, name, active, created_at, updated_at FROM mean_communication WHERE id = $1`

	if err := db.QueryRow(query, id).Scan(&mc.Id, &mc.Name, &mc.Active, &mc.CreatedAt, &mc.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return mc, echo.NewHTTPError(http.StatusNotFound, "Mean not found")
		}
		return mc, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve Mean")
	}

	return mc, nil
}

func InsertComMeans(mcR *models.CreateCommunicationMeanRequest) error {
	db := storage.GetDB()

	query := "INSERT INTO mean_communication (name) VALUES ($1)"

	_, err := db.Exec(query, mcR.Name)
	if err != nil {
		if _, ok := utils.CheckForUserError("unique_type", err); ok {
			return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"errors": echo.Map{"type": "Mean already exists"}})
		}
		return err
	}

	return nil
}

func InsertNegotiation(neg *models.CreateNegotiationRequest) error {
	db := storage.GetDB()

	// 	Name           *string  `json:"Name,omitempty" validate:"required"`
	// Email          *string  `json:"Email,omitempty" validate:"required"`
	// Phone          *string  `json:"Phone,omitempty" validate:"required"`
	// EstimatedValue *float64 `json:"EstimatedValue,omitempty" validate:"required"`
	// BoatName       *string  `json:"BoatName,omitempty" validate:"required"`
	// Qualified      *string  `json:"Qualified,omitempty"`
	// QualifiedType  *string  `json:"QualifiedType,omitempty"`

	query := "INSERT INTO customers (name, email, phone, qualified) VALUES ($1, $2, $3, $4)"

	_, err := db.Exec(query, neg.Name, neg.Email, neg.Phone, neg.Qualified)
	if err != nil {
		// if _, ok := utils.CheckForUserError("unique_type", err); ok {
		// 	return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"errors": echo.Map{"type": "Mean already exists"}})
		// }
		return err
	}

	query = "INSERT INTO so_business (id_customer, boat_name, estimated_value) VALUES ($1, $2, $3)"

	_, err = db.Exec(query, 1, neg.BoatName, neg.EstimatedValue)
	if err != nil {
		// if _, ok := utils.CheckForUserError("unique_type", err); ok {
		// 	return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"errors": echo.Map{"type": "Mean already exists"}})
		// }
		return err
	}

	return nil
}

func GetComMeans(pagenum string, limitPerPage string, name string, active string) ([]models.CommunicationMean, int, error) {
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

	var accs []models.CommunicationMean

	conds := []string{}
	args := []interface{}{}
	paramCount := 1

	if name != "" {
		conds = append(conds, fmt.Sprintf("MC.name ILIKE $%d", paramCount))
		args = append(args, "%"+name+"%")
		paramCount++
	}

	if active != "" {
		conds = append(conds, fmt.Sprintf("MC.active = $%d", paramCount))
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
	SELECT MC.id, MC.name, MC.created_at, MC.updated_at, MC.active
	FROM mean_communication AS MC
	%s
	ORDER BY MC.id, MC.name
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
		var curAcc models.CommunicationMean
		rows.Scan(&curAcc.Id, &curAcc.Name, &curAcc.CreatedAt, &curAcc.UpdatedAt, &curAcc.Active)
		accs = append(accs, curAcc)
	}

	return accs, numRecords, nil
}
