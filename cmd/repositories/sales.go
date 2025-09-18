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

	query := "INSERT INTO customers (id_user, id_mean_communication, name, email, phone, qualified, qualified_type) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"

	var customerID int
	err := db.QueryRow(query, neg.UserId, neg.ComMeanId, neg.Name, neg.Email, neg.Phone, neg.Qualified, neg.QualifiedType).Scan(&customerID)
	if err != nil {
		// if _, ok := utils.CheckForUserError("unique_type", err); ok {
		// 	return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"errors": echo.Map{"type": "Mean already exists"}})
		// }
		return err
	}

	query = "INSERT INTO so_business (id_customer, id_mean_communication, boat_name, estimated_value, qualified, qualified_type) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err = db.Exec(query, customerID, neg.ComMeanId, neg.BoatName, neg.EstimatedValue, neg.Qualified, neg.QualifiedType)
	if err != nil {
		// if _, ok := utils.CheckForUserError("unique_type", err); ok {
		// 	return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"errors": echo.Map{"type": "Mean already exists"}})
		// }
		return err
	}

	return nil
}

func GetCustomersBirthday() ([]models.Customer, int, error) {
	db := storage.GetDB()
	var custs []models.Customer

	query := fmt.Sprintf(`
	SELECT C.id, C.id_user, C.id_mean_communication, U.name AS seller_name, MC.name,
	C.name, C.email, C.phone, C.birthdate, C.pf_pj, 
	C.cpf, C.cnpj, C.cep, C.street, C.neighborhood,
	C.city, C.complement, C.qualified, C.active, C.active_contact

	FROM customers AS C
	INNER JOIN users AS U ON C.id_user = U.id
	INNER JOIN mean_communication AS MC ON C.id_mean_communication = MC.id

	WHERE C.birthdate IS NOT NULL AND
	(EXTRACT(MONTH FROM C.birthdate) = EXTRACT(MONTH FROM CURRENT_DATE) 
    AND EXTRACT(DAY FROM C.birthdate) >= EXTRACT(DAY FROM CURRENT_DATE)
    AND EXTRACT(DAY FROM C.birthdate) <= EXTRACT(DAY FROM CURRENT_DATE + INTERVAL '1 month'))
    	OR
    (EXTRACT(MONTH FROM C.birthdate) = EXTRACT(MONTH FROM CURRENT_DATE + INTERVAL '1 month')
    AND EXTRACT(DAY FROM C.birthdate) <= EXTRACT(DAY FROM CURRENT_DATE + INTERVAL '1 month'))

	ORDER BY EXTRACT(MONTH FROM C.birthdate), EXTRACT(DAY FROM C.birthdate), C.name
	`)

	rows, err := db.Query(query)

	if err != nil {
		if err == sql.ErrNoRows {
			return custs, 0, echo.NewHTTPError(http.StatusNotFound, "Types not found")
		}
		return custs, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve accs"+err.Error())
	}

	queryTotalRecords := fmt.Sprintf(`
	SELECT COUNT(1)
	FROM customers AS C
	INNER JOIN users AS U ON C.id_user = U.id
	INNER JOIN mean_communication AS MC ON C.id_mean_communication = MC.id

	WHERE C.birthdate IS NOT NULL AND
	(EXTRACT(MONTH FROM C.birthdate) = EXTRACT(MONTH FROM CURRENT_DATE) 
    AND EXTRACT(DAY FROM C.birthdate) >= EXTRACT(DAY FROM CURRENT_DATE)
    AND EXTRACT(DAY FROM C.birthdate) <= EXTRACT(DAY FROM CURRENT_DATE + INTERVAL '1 month'))
    	OR
    (EXTRACT(MONTH FROM C.birthdate) = EXTRACT(MONTH FROM CURRENT_DATE + INTERVAL '1 month')
    AND EXTRACT(DAY FROM C.birthdate) <= EXTRACT(DAY FROM CURRENT_DATE + INTERVAL '1 month'))

	ORDER BY EXTRACT(MONTH FROM C.birthdate), EXTRACT(DAY FROM C.birthdate), C.name
	`)

	rowsCount := db.QueryRow(queryTotalRecords)
	numRecords := 0
	rowsCount.Scan(&numRecords)

	for rows.Next() {
		var curC models.Customer
		rows.Scan(&curC.Id, &curC.UserId, &curC.MeanComId, &curC.SellerName, &curC.MeamComName, &curC.Name, &curC.Email, &curC.Phone, &curC.BirthDate, &curC.PfPj, &curC.Cpf, &curC.Cnpj, &curC.Cep, &curC.Street, &curC.Neighborhood, &curC.City, &curC.Complement, &curC.Qualified, &curC.Active, &curC.ActiveContact)
		custs = append(custs, curC)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return custs, numRecords, nil
}

func GetCustomers(pagenum string, limitPerPage string, name string, email string, phone string) ([]models.Customer, int, error) {
	db := storage.GetDB()

	pagenumber, err := strconv.Atoi(pagenum)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve customers (PG1)")
	}
	limit, err := strconv.Atoi(limitPerPage)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve customers (PG2)")
	}

	offset := (pagenumber - 1) * limit

	var custs []models.Customer

	conds := []string{}
	args := []interface{}{}
	paramCount := 1

	if name != "" {
		conds = append(conds, fmt.Sprintf("C.name ILIKE $%d", paramCount))
		args = append(args, "%"+name+"%")
		paramCount++
	}

	if email != "" {
		conds = append(conds, fmt.Sprintf("C.email ILIKE $%d", paramCount))
		args = append(args, "%"+email+"%")
		paramCount++
	}

	if phone != "" {
		conds = append(conds, fmt.Sprintf("C.phone ILIKE $%d", paramCount))
		args = append(args, "%"+phone+"%")
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
	SELECT C.id, C.id_user, C.id_mean_communication, U.name AS seller_name, MC.name,
	C.name, C.email, C.phone, C.birthdate, C.pf_pj, 
	C.cpf, C.cnpj, C.cep, C.street, C.neighborhood,
	C.city, C.complement, C.qualified, C.active, C.active_contact

	FROM customers AS C
	INNER JOIN users AS U ON C.id_user = U.id
	INNER JOIN mean_communication AS MC ON C.id_mean_communication = MC.id
	%s
	ORDER BY C.id, C.name
	LIMIT $%d OFFSET $%d
	`, where, limitArgPos, offsetArgPos)

	rows, err := db.Query(query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return custs, 0, echo.NewHTTPError(http.StatusNotFound, "Types not found")
		}
		return custs, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve accs"+err.Error())
	}

	queryTotalRecords := fmt.Sprintf(`
	SELECT COUNT(1)
	FROM customers AS C
	INNER JOIN users AS U ON C.id_user = U.id
	INNER JOIN mean_communication AS MC ON C.id_mean_communication = MC.id
	%s
	`, where)

	rowsCount := db.QueryRow(queryTotalRecords, args[:len(args)-2]...) // slice to remove the limit and offset args, they are not needed here
	numRecords := 0
	rowsCount.Scan(&numRecords)

	for rows.Next() {
		var curC models.Customer
		rows.Scan(&curC.Id, &curC.UserId, &curC.MeanComId, &curC.SellerName, &curC.MeamComName, &curC.Name, &curC.Email, &curC.Phone, &curC.BirthDate, &curC.PfPj, &curC.Cpf, &curC.Cnpj, &curC.Cep, &curC.Street, &curC.Neighborhood, &curC.City, &curC.Complement, &curC.Qualified, &curC.Active, &curC.ActiveContact)
		custs = append(custs, curC)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return custs, numRecords, nil
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

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return accs, numRecords, nil
}

func GetNegotiations(search string, userId int) ([]models.Negotiation, int, error) {
	db := storage.GetDB()

	var negs []models.Negotiation

	conds := []string{}
	args := []interface{}{}
	paramCount := 1

	if search != "" {
		conds = append(conds, fmt.Sprintf("SB.boat_name ILIKE $%d OR C.name ILIKE $%d", paramCount, paramCount))
		args = append(args, "%"+search+"%")
		paramCount++
	}

	conds = append(conds, fmt.Sprintf("C.id_user = $%d", paramCount))
	args = append(args, userId)
	paramCount++

	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}

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
			C.qualified
	FROM so_business AS SB

	INNER JOIN customers AS C ON SB.id_customer = C.id
	INNER JOIN mean_communication AS MC ON SB.id_mean_communication = MC.id

	%s
	ORDER BY SB.id
	`, where)

	rows, err := db.Query(query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return negs, 0, echo.NewHTTPError(http.StatusNotFound, "Negotiations not found")
		}
		return negs, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve negotiations"+err.Error())
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
			&curNeg.NavigationCity, &curNeg.BoatCapacityNeeded, &curNeg.NewUsed, &curNeg.CabOpen, &curNeg.Stage, &curNeg.Qualified); err != nil {
			return nil, 0, fmt.Errorf("scan error: %w", err)
		}

		negs = append(negs, curNeg)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return negs, numRecords, nil
}
