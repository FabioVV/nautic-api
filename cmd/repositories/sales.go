package repositories

import (
	"database/sql"
	"fmt"
	"nautic/cmd/audit"
	"nautic/cmd/storage"
	"nautic/cmd/utils"
	"nautic/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func UpdateComMean(id int, id_user int, mcR *models.UpdateCommunicationMeanRequest) error {
	db := storage.GetDB()

	mc, err := GetComMean(id, id_user)
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

	// Don't check log insert error for now
	_ = audit.InsertLog(id_user, "/sales/communication-means/:id", "UPDATE", fmt.Sprintf("UPDATE request on communication mean of id %d", id), query)

	return nil
}

func DeactivateComMean(id int, id_user int) error {
	db := storage.GetDB()

	_, err := GetComMean(id, id_user)
	if err != nil {
		return err
	}

	query := `UPDATE mean_communication SET active = 'N' WHERE id = $1`

	_, err = db.Exec(query, id)
	if err != nil {
		return err
	}

	// Don't check log insert error for now
	_ = audit.InsertLog(id_user, "/sales/communication-means/:id", "DELETE", fmt.Sprintf("DELETE request on communication mean of id %d", id), query)

	return nil

}

func GetComMean(id int, id_user int) (models.CommunicationMean, error) {
	db := storage.GetDB()

	var mc models.CommunicationMean
	query := `SELECT id, name, active, created_at, updated_at FROM mean_communication WHERE id = $1`

	if err := db.QueryRow(query, id).Scan(&mc.Id, &mc.Name, &mc.Active, &mc.CreatedAt, &mc.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return mc, echo.NewHTTPError(http.StatusNotFound, "Mean not found")
		}
		return mc, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve Mean")
	}

	// Don't check log insert error for now
	_ = audit.InsertLog(id_user, "/sales/communication-means/:id", "GET", fmt.Sprintf("GET request on communication mean of id %d", id), query)

	return mc, nil
}

func InsertSalesOrderUsingBusinessHistory(idBusinessHistory int, id_user int) error {
	db := storage.GetDB()

	query := "INSERT INTO sales_orders (id_business_history) VALUES ($1) RETURNING id"

	var soID int
	err := db.QueryRow(query, idBusinessHistory).Scan(&soID)
	if err != nil {
		return err
	}

	query = "UPDATE business_histories SET id_sales_order = $1 WHERE id = $2"

	_, err = db.Exec(query, soID, idBusinessHistory)
	if err != nil {
		return err
	}

	// Don't check log insert error for now
	_ = audit.InsertLog(id_user, "/orders/communication-means/:id", "GET", fmt.Sprintf("GET request on communication mean of id %d", idBusinessHistory), query)

	return nil
}

func InsertComMeans(id_user int, mcR *models.CreateCommunicationMeanRequest) error {
	db := storage.GetDB()

	query := "INSERT INTO mean_communication (name) VALUES ($1)"

	_, err := db.Exec(query, mcR.Name)
	if err != nil {
		if _, ok := utils.CheckForError("unique_type", err); ok {
			return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"errors": echo.Map{"type": "Mean already exists"}})
		}
		return err
	}

	// Don't check log insert error for now
	_ = audit.InsertLog(id_user, "/sales/communication-means", "POST", fmt.Sprintf("POST request on communication mean"), query)

	return nil
}

func InsertNegotiation(id_user int, neg *models.CreateNegotiationRequest) error {
	db := storage.GetDB()

	query := "INSERT INTO customers (id_user, id_mean_communication, name, email, phone, qualified, qualified_type, boat_alert) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"

	var customerID int
	err := db.QueryRow(query, neg.UserId, neg.ComMeanId, neg.Name, neg.Email, neg.Phone, neg.Qualified, neg.QualifiedType, neg.BoatName).Scan(&customerID)
	if err != nil {
		// if _, ok := utils.CheckForError("unique_type", err); ok {
		// 	return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"errors": echo.Map{"type": "Mean already exists"}})
		// }
		return err
	}

	query = "INSERT INTO so_business (id_customer, id_mean_communication, boat_name, estimated_value, qualified, qualified_type) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err = db.Exec(query, customerID, neg.ComMeanId, neg.BoatName, neg.EstimatedValue, neg.Qualified, neg.QualifiedType)
	if err != nil {
		return err
	}

	// Don't check log insert error for now
	_ = audit.InsertLog(id_user, "/sales/negotiations", "POST", fmt.Sprintf("POST request on negotiations"), query)

	return nil
}

func CreateNegotiationHistory(id int, neg *models.CreateNegotiationHistoryRequest) error {
	db := storage.GetDB()

	if id == 0 {
		query := "INSERT INTO business_histories (id_user, id_customer, description, stage, id_mean_communication, id_business) VALUES ($1, $2, $3, $4, $5, $6)"

		_, err := db.Exec(query, neg.UserId, neg.CustomerId, neg.Description, neg.Stage, neg.ComMeanId, nil)
		if err != nil {
			// if _, ok := utils.CheckForError("unique_type", err); ok {
			// 	return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"errors": echo.Map{"type": "Mean already exists"}})
			// }
			return err
		}

		return nil
	}

	query := "INSERT INTO business_histories (id_user, id_customer, description, stage, id_mean_communication, id_business) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := db.Exec(query, neg.UserId, neg.CustomerId, neg.Description, neg.Stage, neg.ComMeanId, id)
	if err != nil {
		// if _, ok := utils.CheckForError("unique_type", err); ok {
		// 	return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"errors": echo.Map{"type": "Mean already exists"}})
		// }
		return err
	}

	// Don't check log insert error for now
	_ = audit.InsertLog(id, "/sales/negotiation/:id/history", "POST", fmt.Sprintf("POST request on negotiations histories"), query)

	return nil
}

func GetCustomersBirthday(id_user int) ([]models.Customer, int, error) {
	db := storage.GetDB()
	var custs []models.Customer

	query := `
	SELECT C.id, C.id_user, C.id_mean_communication, U.name AS seller_name, MC.name,
	C.name, C.email, C.phone, C.birthdate, C.pf_pj, 
	C.cpf, C.cnpj, C.cep, C.street, C.neighborhood,
	C.city, C.complement, C.qualified, C.active, C.active_contact

	FROM customers AS C
	INNER JOIN users AS U ON C.id_user = U.id
	INNER JOIN mean_communication AS MC ON C.id_mean_communication = MC.id

	WHERE C.birthdate IS NOT NULL 
	AND (
        make_date(
            EXTRACT(YEAR FROM CURRENT_DATE)::int,
            EXTRACT(MONTH FROM C.birthdate)::int,
            EXTRACT(DAY   FROM C.birthdate)::int
        ) BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '30 days'
        OR
        make_date(
            (EXTRACT(YEAR FROM CURRENT_DATE)::int) + 1,
            EXTRACT(MONTH FROM C.birthdate)::int,
            EXTRACT(DAY   FROM C.birthdate)::int
        ) BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '30 days'
      )
	ORDER BY EXTRACT(MONTH FROM C.birthdate), EXTRACT(DAY   FROM C.birthdate), C.name
	`

	rows, err := db.Query(query)

	if err != nil {
		if err == sql.ErrNoRows {
			return custs, 0, echo.NewHTTPError(http.StatusNotFound, "Types not found")
		}
		return custs, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve accs"+err.Error())
	}

	queryTotalRecords := `
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
	`

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

	// Don't check log insert error for now
	_ = audit.InsertLog(id_user, "/customers/birthdays", "GET", fmt.Sprintf("GET request on customer birthdays"), query)

	return custs, numRecords, nil
}

func GetCustomers(id_user int, pagenum string, limitPerPage string, name string, email string, phone string, boat string) ([]models.Customer, int, error) {
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

	if boat != "" {
		conds = append(conds, fmt.Sprintf("C.boat_alert ILIKE $%d", paramCount))
		args = append(args, "%"+boat+"%")
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

	// Don't check log insert error for now
	_ = audit.InsertLog(id_user, "/customers", "GET", fmt.Sprintf("GET request on customers"), query)

	return custs, numRecords, nil
}

func GetComMeans(id_user int, pagenum string, limitPerPage string, name string, active string) ([]models.CommunicationMean, int, error) {
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

	// Don't check log insert error for now
	_ = audit.InsertLog(id_user, "/sales/communication-means", "GET", fmt.Sprintf("GET request on communication means"), query)

	return accs, numRecords, nil
}

func GetNegotiationHistory(id_business int, id_user int) ([]models.NegotiationHistory, int, error) { // aa
	db := storage.GetDB()

	var negsh []models.NegotiationHistory

	// conds := []string{}
	// args := []interface{}{}

	// where := ""
	// if len(conds) > 0 {
	// 	where = "WHERE " + strings.Join(conds, " AND ")
	// }

	query := `
	SELECT BIH.id, BIH.id_user, BIH.id_customer, BIH.id_mean_communication, 
	BIH.description, BIH.stage, BIH.created_at,
	C.name, MC.name, BIH.id_business, SB.negotiation_active, BIH.id_sales_order,
	EXISTS (
        SELECT 1
        FROM sales_orders as SO
        WHERE SO.id_business_history = BIH.id 
    ) AS has_sales_order,
	(SELECT status FROM sales_orders WHERE id = BIH.id_sales_order) AS sales_order_canceled,
	(SELECT done FROM sales_orders WHERE id = BIH.id_sales_order) AS sales_order_finished

	FROM business_histories AS BIH

	INNER JOIN customers AS C ON BIH.id_customer = C.id
	INNER JOIN mean_communication AS MC ON BIH.id_mean_communication = MC.id
	INNER JOIN so_business AS SB ON BIH.id_business = SB.id AND SB.id = $1

	WHERE BIH.id_user = $2

	ORDER BY BIH.id DESC
	`

	rows, err := db.Query(query, id_business, id_user)

	if err != nil {
		if err == sql.ErrNoRows {
			return negsh, 0, echo.NewHTTPError(http.StatusNotFound, "Negotiations not found")
		}
		return negsh, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve negotiations")
	}

	queryTotalRecords := `
	SELECT COUNT(1)

	FROM business_histories AS BIH

	INNER JOIN customers AS C ON BIH.id_customer = C.id
	INNER JOIN mean_communication AS MC ON BIH.id_mean_communication = MC.id
	INNER JOIN so_business AS SB ON BIH.id_business = SB.id AND SB.id = $1

	WHERE BIH.id_user = $2


	ORDER BY BIH.id DESC
	`
	//println(queryTotalRecords)

	rowsCount := db.QueryRow(queryTotalRecords, id_business, id_user)
	numRecords := 0
	rowsCount.Scan(&numRecords)

	for rows.Next() {
		var curNegH models.NegotiationHistory

		if err := rows.Scan(&curNegH.Id, &curNegH.UserId, &curNegH.CustomerId, &curNegH.ComMeanId,
			&curNegH.Description, &curNegH.Stage, &curNegH.DateCreated, &curNegH.CustomerName,
			&curNegH.MeamComName, &curNegH.BusinessId,
			&curNegH.NegotiationActive, &curNegH.SalesOrderId,
			&curNegH.HasSalesOrder, &curNegH.SalesOrderCanceled, &curNegH.SalesOrderFinished); err != nil {
			return nil, 0, fmt.Errorf("scan error: %w", err)
		}

		negsh = append(negsh, curNegH)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	// Don't check log insert error for now
	_ = audit.InsertLog(id_user, "/negotiations/:id/history", "GET", fmt.Sprintf("GET request on negotiation history of business %d", id_business), query)

	return negsh, numRecords, nil
}

func GetNegotiationsAlerts(userId int) ([]models.NegotiationAlert, error) {
	db := storage.GetDB()

	var negs []models.NegotiationAlert

	conds := []string{}
	args := []any{}
	paramCount := 1

	conds = append(conds, fmt.Sprintf("C.id_user = $%d AND NA.date = DATE(CURRENT_TIMESTAMP)", paramCount))
	args = append(args, userId)
	paramCount++

	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}

	query := fmt.Sprintf(`
	SELECT NA.id, NA.id_business, C.name, C.phone, NA.motive, NA.date
	FROM negotiations_alerts AS NA

	INNER JOIN so_business AS SB ON NA.id_business = SB.id
	INNER JOIN customers AS C ON SB.id_customer = C.id
	INNER JOIN users AS U ON C.id_user = U.id

	%s
	ORDER BY NA.id
	`, where)

	rows, err := db.Query(query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return negs, echo.NewHTTPError(http.StatusNotFound, "Negotiation alerts not found")
		}
		return negs, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve negotiation alerts")
	}

	queryTotalRecords := fmt.Sprintf(`
	SELECT COUNT(1)
	FROM negotiations_alerts AS NA

	INNER JOIN so_business AS SB ON NA.id_business = SB.id
	INNER JOIN customers AS C ON SB.id_customer = C.id
	INNER JOIN users AS U ON C.id_user = U.id

	%s
	`, where)
	//println(queryTotalRecords)

	rowsCount := db.QueryRow(queryTotalRecords, args...)
	numRecords := 0
	rowsCount.Scan(&numRecords)

	for rows.Next() {
		var curNeg models.NegotiationAlert

		if err := rows.Scan(&curNeg.Id, &curNeg.BusinessId, &curNeg.CustomerName, &curNeg.CustomerPhone, &curNeg.Motive, &curNeg.Date); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		negs = append(negs, curNeg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	// Don't check log insert error for now
	_ = audit.InsertLog(userId, "/negotiations/alerts", "GET", fmt.Sprintf("GET request on negotiation alerts"), query)

	return negs, nil
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

	conds = append(conds, fmt.Sprintf("SB.negotiation_active = $%d", paramCount))
	args = append(args, "Y")
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
			C.qualified,
			(SB.created_at < now() - interval '24 hours') AS has_passed_24hrs,
			-- score calculation
			(
			-- base 0
			COALESCE(
			-- has boat: +100
			CASE WHEN SB.has_boat = 'Y' OR SB.has_boat = 'S' THEN 200 ELSE 0 END
			-- suspect of fraud: -1000
			+ CASE WHEN C.suspect_of_fraud = 'Y' THEN -1000 ELSE 0 END
			-- willing to spend a lot: +100 (assume estimated_value threshold)
			+ CASE WHEN SB.estimated_value IS NOT NULL AND SB.estimated_value >= 300000 THEN 200 ELSE 0 END
			-- slow to return contact: -100 (assume last history older than 7 days)
			+ CASE
			WHEN bh.last_history_at IS NULL THEN -100
			WHEN now() - bh.last_history_at > interval '7 days' THEN -500
			ELSE 0
			END
			+ CASE
			WHEN now() - SB.stage_last_updated_at > interval '7 days' THEN -500
			WHEN now() - SB.stage_last_updated_at > interval '5 days' THEN -300
			WHEN now() - SB.stage_last_updated_at > interval '3 days' THEN -100
			WHEN now() - SB.stage_last_updated_at > interval '1 days' OR now() - SB.stage_last_updated_at > interval '2 days' THEN 100
			ELSE 0
			END
			+ CASE WHEN C.qualified = 'Y' THEN 300 ELSE -100 END
			+ CASE WHEN C.qualified = 'Y' AND C.qualified_type = 'A' THEN 300 ELSE 0 END
			+ CASE WHEN C.qualified = 'Y' AND C.qualified_type = 'B' THEN 200 ELSE 0 END
			+ CASE WHEN C.qualified = 'Y' AND C.qualified_type = 'C' THEN 100 ELSE 0 END

			+ CASE WHEN SB.stage = 1 THEN 100 ELSE 0 END
			+ CASE WHEN SB.stage = 2 THEN 100 ELSE 0 END
			+ CASE WHEN SB.stage = 3 THEN 200 ELSE 0 END
			+ CASE WHEN SB.stage = 4 THEN 300 ELSE 0 END

			,0)
			) AS customer_score


	FROM so_business AS SB

	INNER JOIN customers AS C ON SB.id_customer = C.id
	INNER JOIN mean_communication AS MC ON SB.id_mean_communication = MC.id
	LEFT JOIN (
		-- last business_histories.created_at per business (id_business links to so_business.id)
		SELECT id_business,
		max(created_at) AS last_history_at
		FROM public.business_histories
		GROUP BY id_business
	) bh ON bh.id_business = SB.id

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
			&curNeg.NavigationCity, &curNeg.BoatCapacityNeeded, &curNeg.NewUsed, &curNeg.CabOpen, &curNeg.Stage,
			&curNeg.Qualified, &curNeg.HasPassed24Hrs, &curNeg.CustomerScore); err != nil {
			return nil, 0, fmt.Errorf("scan error: %w", err)
		}

		negs = append(negs, curNeg)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	// Don't check log insert error for now
	_ = audit.InsertLog(userId, "/negotiations", "GET", fmt.Sprintf("GET request on negotiations"), query)

	return negs, numRecords, nil
}

func GetCustomerNegotiationsHistories(id_user int, customerId int) ([]models.NegotiationHistory, int, error) {
	db := storage.GetDB()

	var negsh []models.NegotiationHistory

	// conds := []string{}
	// args := []interface{}{}

	// where := ""
	// if len(conds) > 0 {
	// 	where = "WHERE " + strings.Join(conds, " AND ")
	// }

	query := `
	SELECT BIH.id, BIH.id_user, BIH.id_customer, BIH.id_mean_communication, 
	BIH.description, BIH.stage, BIH.created_at, BIH.id_business, C.name, MC.name, SB.negotiation_active, BIH.id_sales_order,
	EXISTS (
        SELECT 1
        FROM sales_orders as SO
        WHERE SO.id_business_history = BIH.id 
    ) AS has_sales_order

	FROM business_histories AS BIH

	INNER JOIN customers AS C ON BIH.id_customer = C.id
	INNER JOIN mean_communication AS MC ON BIH.id_mean_communication = MC.id
	LEFT JOIN so_business AS SB ON BIH.id_business = SB.id

	WHERE C.id = $1

	ORDER BY BIH.id DESC
	`

	rows, err := db.Query(query, customerId)

	if err != nil {
		if err == sql.ErrNoRows {
			return negsh, 0, echo.NewHTTPError(http.StatusNotFound, "Negotiations not found")
		}
		return negsh, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve negotiations")
	}

	queryTotalRecords := `
	SELECT COUNT(1)

	FROM business_histories AS BIH

	INNER JOIN customers AS C ON BIH.id_customer = C.id
	INNER JOIN mean_communication AS MC ON BIH.id_mean_communication = MC.id
	INNER JOIN so_business AS SB ON BIH.id_business = SB.id

	WHERE C.id = $1


	ORDER BY BIH.id
	`
	//println(queryTotalRecords)

	rowsCount := db.QueryRow(queryTotalRecords, customerId)
	numRecords := 0
	rowsCount.Scan(&numRecords)

	for rows.Next() {
		var curNegH models.NegotiationHistory

		// 		SELECT BIH.id, BIH.id_user, BIH.id_customer, BIH.id_mean_communication,
		// BIH.description, BIH.stage, BIH.created_at,
		// C.name, MC.name

		if err := rows.Scan(&curNegH.Id, &curNegH.UserId, &curNegH.CustomerId, &curNegH.ComMeanId,
			&curNegH.Description, &curNegH.Stage, &curNegH.DateCreated, &curNegH.BusinessId, &curNegH.CustomerName, &curNegH.MeamComName, &curNegH.NegotiationActive, &curNegH.SalesOrderId, &curNegH.HasSalesOrder); err != nil {
			return nil, 0, fmt.Errorf("scan error: %w", err)
		}

		negsh = append(negsh, curNegH)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	// Don't check log insert error for now
	_ = audit.InsertLog(id_user, "/customers/:id/negotiations/history", "GET", fmt.Sprintf("GET request on negotiations history of customer of id %d", customerId), query)

	return negsh, numRecords, nil
}

func GetCustomer(id int, id_user int) (models.Customer, error) {
	db := storage.GetDB()

	var curC models.Customer
	query := `SELECT C.id, C.id_user, C.id_mean_communication, U.name AS seller_name, MC.name,
	C.name, C.email, C.phone, C.birthdate, C.pf_pj, 
	C.cpf, C.cnpj, C.cep, C.street, C.neighborhood,
	C.city, C.state, C.complement, C.qualified, C.active, C.active_contact, C.suspect_of_fraud,
	SB_LAST.estimated_value, SB_LAST.customer_city, SB_LAST.customer_navigation_city, SB_LAST.boat_capacity_needed,
	SB_LAST.new_used, SB_LAST.cab_open, SB_LAST.boat_length_min, SB_LAST.boat_length_max

	FROM customers AS C
	INNER JOIN users AS U ON C.id_user = U.id
	INNER JOIN mean_communication AS MC ON C.id_mean_communication = MC.id
	LEFT JOIN LATERAL (
		SELECT *
		FROM so_business sb
		WHERE sb.id_customer = C.id
		ORDER BY sb.created_at DESC    -- or ORDER BY sb.id DESC
		LIMIT 1
	) AS SB_LAST ON true

	WHERE C.id = $1
	ORDER BY C.id, C.name
	`

	if err := db.QueryRow(query, id).Scan(&curC.Id, &curC.UserId, &curC.MeanComId,
		&curC.SellerName, &curC.MeamComName, &curC.Name,
		&curC.Email, &curC.Phone, &curC.BirthDate, &curC.PfPj,
		&curC.Cpf, &curC.Cnpj, &curC.Cep, &curC.Street, &curC.Neighborhood,
		&curC.City, &curC.State, &curC.Complement, &curC.Qualified, &curC.Active,
		&curC.ActiveContact, &curC.SuspectOfFraud, &curC.EstimatedValue, &curC.CustomerCity, &curC.NavigationCity,
		&curC.BoatCapacity, &curC.NewUsed, &curC.CabinatedOpen, &curC.MinPesBoat, &curC.MaxPesBoat); err != nil {
		if err == sql.ErrNoRows {
			return curC, echo.NewHTTPError(http.StatusNotFound, "Negotiation not found")
		}
		return curC, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve negotiation")
	}

	// Don't check log insert error for now
	_ = audit.InsertLog(id_user, "/customers", "GET", fmt.Sprintf("GET request on customer of id %d", id), query)
	return curC, nil
}

func GetNegotiation(id_user int, id int) (models.Negotiation, error) {
	db := storage.GetDB()

	var curNeg models.Negotiation
	query := `SELECT SB.id, 
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
			SB.has_boat, 
			SB.has_boat_which, 
			SB.boat_length_min, 
			SB.boat_length_max, 
			C.qualified,
			C.qualified_type
			
	FROM so_business AS SB

	INNER JOIN customers AS C ON SB.id_customer = C.id
	INNER JOIN mean_communication AS MC ON SB.id_mean_communication = MC.id
	WHERE SB.id = $1`

	if err := db.QueryRow(query, id).Scan(&curNeg.Id, &curNeg.CustomerId, &curNeg.MeanComId,
		&curNeg.Name, &curNeg.Email, &curNeg.Phone, &curNeg.MeamComName,
		&curNeg.BoatName, &curNeg.EstimatedValue, &curNeg.MaxEstimatedValue, &curNeg.City,
		&curNeg.NavigationCity, &curNeg.BoatCapacityNeeded, &curNeg.NewUsed, &curNeg.CabOpen,
		&curNeg.Stage, &curNeg.HasBoat, &curNeg.HasBoatWhich, &curNeg.MinPesBoat, &curNeg.MaxPesBoat,
		&curNeg.Qualified, &curNeg.QualifiedType); err != nil {
		if err == sql.ErrNoRows {
			return curNeg, echo.NewHTTPError(http.StatusNotFound, "Negotiation not found")
		}
		return curNeg, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve negotiation"+err.Error())
	}

	// Don't check log insert error for now
	_ = audit.InsertLog(id_user, "/negotiations", "GET", fmt.Sprintf("GET request on negotiation of id %d", id), query)

	return curNeg, nil
}

func UpdateCustomer(id_user int, id int, cT *models.CustomerRequest) error {
	db := storage.GetDB()

	// _, err := GetNegotiation(id)
	// if err != nil {
	// 	return err
	// }

	// if accTg.Active == "N" {
	// 	return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"errors": echo.Map{"accessory": "Engine must bet active to update it"}})
	// }

	query := `UPDATE customers SET `
	params := []interface{}{}
	paramCount := 0

	if cT.Name != nil {
		paramCount++
		query += fmt.Sprintf("name = $%d, ", paramCount)
		params = append(params, *&cT.Name)
	}

	if cT.Email != nil {
		paramCount++
		query += fmt.Sprintf("email = $%d, ", paramCount)
		params = append(params, *&cT.Email)
	}

	if cT.Phone != nil {
		paramCount++
		query += fmt.Sprintf("phone = $%d, ", paramCount)
		params = append(params, *&cT.Phone)
	}

	if cT.Cep != nil {
		paramCount++
		query += fmt.Sprintf("cep = $%d, ", paramCount)
		params = append(params, *&cT.Cep)
	}

	if cT.Street != nil {
		paramCount++
		query += fmt.Sprintf("street = $%d, ", paramCount)
		params = append(params, *&cT.Street)
	}

	if cT.Neighborhood != nil {
		paramCount++
		query += fmt.Sprintf("neighborhood = $%d, ", paramCount)
		params = append(params, *&cT.Neighborhood)
	}

	if cT.Complement != nil {
		paramCount++
		query += fmt.Sprintf("complement = $%d, ", paramCount)
		params = append(params, *&cT.Complement)
	}

	if cT.State != nil {
		paramCount++
		query += fmt.Sprintf("state = $%d, ", paramCount)
		params = append(params, *&cT.State)
	}

	if cT.City != nil {
		paramCount++
		query += fmt.Sprintf("city = $%d, ", paramCount)
		params = append(params, *&cT.City)
	}

	if cT.HasBoat != nil {
		paramCount++
		query += fmt.Sprintf("has_boat = $%d, ", paramCount)
		params = append(params, *&cT.HasBoat)
	}

	if cT.WhichBoat != nil {
		paramCount++
		query += fmt.Sprintf("has_boat_which = $%d, ", paramCount)
		params = append(params, *&cT.WhichBoat)
	}

	if cT.BirthDay != nil {
		paramCount++
		query += fmt.Sprintf("birthdate = $%d, ", paramCount)
		params = append(params, *&cT.BirthDay)
	}

	if cT.PfPj != nil {
		paramCount++
		query += fmt.Sprintf("pf_pj = $%d, ", paramCount)
		params = append(params, *&cT.PfPj)
	}

	if cT.Cpf != nil {
		paramCount++
		query += fmt.Sprintf("cpf = $%d, ", paramCount)
		params = append(params, *&cT.Cpf)
	}

	if cT.Cnpj != nil {
		paramCount++
		query += fmt.Sprintf("cnpj = $%d, ", paramCount)
		params = append(params, *&cT.Cnpj)
	}

	if cT.SuspectOfFraud != nil {
		paramCount++
		query += fmt.Sprintf("suspect_of_fraud = $%d, ", paramCount)
		params = append(params, *&cT.SuspectOfFraud)
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

	// Don't check log insert error for now
	_ = audit.InsertLog(id_user, "/customers/:id", "UPDATE", fmt.Sprintf("UPDATE request on customer of id %d", id), query)

	return nil
}

func ReactivateNegotiation(id_user int, id int) error {
	db := storage.GetDB()
	query := `UPDATE so_business SET negotiation_active = 'Y', updated_at = NOW(), stage = 1 WHERE id = $1`

	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	// Don't check log insert error for now
	_ = audit.InsertLog(id_user, "/sales/negotiations/:id/reactivate", "UPDATE", fmt.Sprintf("UPDATE request on negotiation of id %d", id), query)

	return nil
}

func LostNegotiation(id_user int, id int, negT *models.LostNegotiation) error {
	db := storage.GetDB()
	query := `UPDATE so_business SET negotiation_active = $1, updated_at = NOW() WHERE id = $2`

	_, err := db.Exec(query, "N", id)
	if err != nil {
		return err
	}

	query = `INSERT INTO lost_negotiations (id_business, motive, customer_got_another_boat, which_boat, our_boat_offered) VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(query, id, negT.Motive, negT.CustomerGotAnotherBoat, negT.WhichBoat, negT.OurBoatOffered)
	if err != nil {
		return err
	}

	query = `INSERT INTO negotiations_alerts (id_business, motive, date) VALUES ($1, $2, $3)`
	_, err = db.Exec(query, id, negT.DateAlertMotive, negT.DateAlert)
	if err != nil {
		return err
	}

	// Don't check log insert error for now
	_ = audit.InsertLog(id_user, "/sales/negotiations/:id/deactivate", "UPDATE", fmt.Sprintf("UPDATE request on negotiation of id %d", id), query)

	return nil
}

func UpdateNegotiation(id_user int, id int, negT *models.CreateNegotiationRequest) error {
	db := storage.GetDB()

	// _, err := GetNegotiation(id)
	// if err != nil {
	// 	return err
	// }

	// if accTg.Active == "N" {
	// 	return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"errors": echo.Map{"accessory": "Engine must bet active to update it"}})
	// }

	query := `UPDATE so_business SET `
	params := []interface{}{}
	paramCount := 0

	if negT.EstimatedValue != nil {
		paramCount++
		query += fmt.Sprintf("estimated_value = $%d, ", paramCount)
		params = append(params, *&negT.EstimatedValue)
	}

	if negT.BoatName != nil {
		paramCount++
		query += fmt.Sprintf("boat_name = $%d, ", paramCount)
		params = append(params, *&negT.BoatName)
	}

	if negT.Qualified != nil {
		paramCount++
		query += fmt.Sprintf("qualified = $%d, ", paramCount)
		params = append(params, *&negT.Qualified)
	}

	if negT.QualifiedType != nil {
		paramCount++
		query += fmt.Sprintf("qualified_type = $%d, ", paramCount)
		params = append(params, *&negT.QualifiedType)
	}

	if negT.City != nil {
		paramCount++
		query += fmt.Sprintf("customer_city = $%d, ", paramCount)
		params = append(params, *&negT.City)
	}

	if negT.NavigationCity != nil {
		paramCount++
		query += fmt.Sprintf("customer_navigation_city = $%d, ", paramCount)
		params = append(params, *&negT.NavigationCity)
	}

	if negT.BoatCapacity != nil {
		paramCount++
		query += fmt.Sprintf("boat_capacity_needed = $%d, ", paramCount)
		params = append(params, *&negT.BoatCapacity)
	}

	if negT.CabinatedOpen != nil {
		paramCount++
		query += fmt.Sprintf("cab_open = $%d, ", paramCount)
		params = append(params, *&negT.CabinatedOpen)
	}

	if negT.ComMeanId != nil {
		paramCount++
		query += fmt.Sprintf("id_mean_communication = $%d, ", paramCount)
		params = append(params, *&negT.ComMeanId)
	}

	if negT.NewUsed != nil {
		paramCount++
		query += fmt.Sprintf("new_used = $%d, ", paramCount)
		params = append(params, *&negT.NewUsed)
	}

	if negT.HasBoat != nil {
		paramCount++
		query += fmt.Sprintf("has_boat = $%d, ", paramCount)
		params = append(params, *&negT.HasBoat)
	}

	if negT.MinPesBoat != nil {
		paramCount++
		query += fmt.Sprintf("boat_length_min = $%d, ", paramCount)
		params = append(params, *&negT.MinPesBoat)
	}

	if negT.MaxPesBoat != nil {
		paramCount++
		query += fmt.Sprintf("boat_length_max = $%d, ", paramCount)
		params = append(params, *&negT.MaxPesBoat)
	}

	if negT.WhichBoat != nil {
		paramCount++
		query += fmt.Sprintf("has_boat_which = $%d, ", paramCount)
		params = append(params, *&negT.WhichBoat)
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

	// Don't check log insert error for now
	_ = audit.InsertLog(id_user, "/sales/negotiations/:id", "UPDATE", fmt.Sprintf("UPDATE request on negotiation of id %d", id), query)

	return nil
}

func UpgradeQuoteToSalesOrder(id_user int, id int) error {
	db := storage.GetDB()

	salesOr, err := GetSalesOrder(id)
	if err != nil {
		return err
	}

	if *salesOr.StatusType == "Orçamento cancelado" || *salesOr.StatusType == "Pedido cancelado" || *salesOr.StatusType == "Pedido concluído" {
		return echo.NewHTTPError(http.StatusInternalServerError, "Sales order cant be updated")
	}

	var status string = ""
	query := `SELECT status FROM sales_orders WHERE id = $1`
	db.QueryRow(query, id).Scan(&status)

	query = `UPDATE sales_orders SET status = $1 WHERE id = $2`

	if status == "NQ" {
		_, err := db.Exec(query, "NO", id)
		if err != nil {
			return err
		}

		// Don't check log insert error for now
		_ = audit.InsertLog(id_user, "/sales/orders/:id/upgrade-quote", "UPDATE", fmt.Sprintf("UPDATE request on quote of id %d", id), query)

		return nil
	}

	return echo.NewHTTPError(http.StatusBadRequest, "Sales order cannot be upgraded to order")
}

func CancelSalesOrder(id_user int, id int) error {
	db := storage.GetDB()

	var status string = ""
	query := `SELECT status FROM sales_orders WHERE id = $1`
	db.QueryRow(query, id).Scan(&status)

	query = `UPDATE sales_orders SET status = $1 WHERE id = $2`
	switch status {
	case "NQ":
		_, err := db.Exec(query, "QC", id)
		if err != nil {
			return err
		}
	case "NO ":
		_, err := db.Exec(query, "OC", id)
		if err != nil {
			return err
		}
	default:
		_, err := db.Exec(query, "QC", id)
		if err != nil {
			return err
		}
	}

	// Don't check log insert error for now
	_ = audit.InsertLog(id_user, "/orders/:id", "DELETE", fmt.Sprintf("DELETE request on sales order of id %d", id), query)

	return nil
}

func UpdateEngineSalesOrder(id_user int, id int, id_engine int) error {
	db := storage.GetDB()

	engine, err := GetEngine(id_engine)
	if err != nil {
		return err
	}

	salesOr, err := GetSalesOrder(id)
	if err != nil {
		return err
	}

	if *salesOr.StatusType == "Orçamento cancelado" || *salesOr.StatusType == "Pedido cancelado" || *salesOr.StatusType == "Pedido concluído" {
		return echo.NewHTTPError(http.StatusInternalServerError, "Sales order cant be updated")
	}

	var id_check int = 0
	query := `SELECT id FROM sales_orders_itens WHERE id_engine = $1 AND id_sales_order = $2`
	db.QueryRow(query, id_engine, id).Scan(&id_check)

	if id_check == 0 {
		query := `INSERT INTO sales_orders_itens(id_engine, id_sales_order, qty, unit_price) VALUES ($1, $2, $3, $4)`

		_, err := db.Exec(query, id_engine, id, 1, engine.PriceSell)
		if err != nil {
			return err
		}
	} else {
		query := `UPDATE sales_orders_itens SET id_engine = $1, unit_price = $2 WHERE id_sales_order = $3`

		_, err := db.Exec(query, id_engine, engine.PriceSell, id)
		if err != nil {
			return err
		}
	}

	// Don't check log insert error for now
	_ = audit.InsertLog(id_user, "/orders/:id/engine/:id_engine", "UPDATE", fmt.Sprintf("UPDATE request on sales order of id %d", id), query)

	return nil
}

func UpdateBoatSalesOrder(id_user int, id int, id_boat int) error {
	db := storage.GetDB()

	boat, err := GetBoat(id_user, id_boat)
	if err != nil {
		return err
	}

	salesOr, err := GetSalesOrder(id)
	if err != nil {
		return err
	}

	if *salesOr.StatusType == "Orçamento cancelado" || *salesOr.StatusType == "Pedido cancelado" || *salesOr.StatusType == "Pedido concluído" {
		return echo.NewHTTPError(http.StatusInternalServerError, "Sales order cant be updated")
	}

	var id_check int = 0
	query := `SELECT id FROM sales_orders_itens WHERE id_boat = $1 AND id_sales_order = $2`
	db.QueryRow(query, id_boat, id).Scan(&id_check)

	if id_check == 0 {
		query := `INSERT INTO sales_orders_itens(id_boat, id_sales_order, qty, unit_price) VALUES ($1, $2, $3, $4)`

		_, err := db.Exec(query, id_boat, id, 1, boat.PriceSell)
		if err != nil {
			return err
		}
	} else {
		query := `UPDATE sales_orders_itens SET id_boat = $1, unit_price = $2 WHERE id_sales_order = $3`

		_, err := db.Exec(query, id_boat, boat.PriceSell, id)
		if err != nil {
			return err
		}
	}

	// Don't check log insert error for now
	_ = audit.InsertLog(id_user, "/orders/:id/boat/:id_boat", "UPDATE", fmt.Sprintf("UPDATE request on sales order of id %d, change boat", id), query)

	return nil
}

func RemoveAccessorySalesOrder(id_user int, id int, id_accessory int) error {
	db := storage.GetDB()

	// _, err := GetAccessory(id_accessory)
	// if err != nil {
	// 	return err
	// }

	query := `DELETE FROM sales_orders_itens WHERE id_accessory = $1 AND id_sales_order = $2`

	_, err := db.Exec(query, id_accessory, id)
	if err != nil {
		return err
	}

	// Don't check log insert error for now
	_ = audit.InsertLog(id_user, "/orders/:id/accessory/:id_accessory", "DELETE", fmt.Sprintf("DELETE request on sales order of id %d, remove accessory", id), query)

	return nil
}

func UpdateAccessoryQtySalesOrder(id_user int, id int, id_accessory int, update *models.UpdateSalesOrderItemQty) error {
	db := storage.GetDB()

	_, err := GetAccessory(id_accessory, id_user)
	if err != nil {
		return err
	}

	salesOr, err := GetSalesOrder(id)
	if err != nil {
		return err
	}

	if *salesOr.StatusType == "Orçamento cancelado" || *salesOr.StatusType == "Pedido cancelado" || *salesOr.StatusType == "Pedido concluído" {
		return echo.NewHTTPError(http.StatusInternalServerError, "Sales order cant be updated")
	}

	var id_check int = 0
	query := `SELECT id FROM sales_orders_itens WHERE id_accessory = $1 AND id_sales_order = $2`
	db.QueryRow(query, id_accessory, id).Scan(&id_check)

	if id_check != 0 {
		query := `UPDATE sales_orders_itens SET qty = $1 WHERE id_accessory = $2 AND id_sales_order = $3`

		_, err := db.Exec(query, update.Qty, id_accessory, id)
		if err != nil {
			return err
		}
	} else {
		return echo.NewHTTPError(http.StatusInternalServerError, "Accessory not found on boat")
	}

	// Don't check log insert error for now
	_ = audit.InsertLog(id_user, "/orders/:id/accessory/:id_accessory/change-qty", "UPDATE", fmt.Sprintf("UPDATE request on sales order of id %d, change accessory quantity", id), query)

	return nil
}

func UpdateAccessorySalesOrder(id_user int, id int, id_accessory int) error {
	db := storage.GetDB()

	acc, err := GetAccessory(id_accessory, id_user)
	if err != nil {
		return err
	}

	salesOr, err := GetSalesOrder(id)
	if err != nil {
		return err
	}

	if *salesOr.StatusType == "Orçamento cancelado" || *salesOr.StatusType == "Pedido cancelado" || *salesOr.StatusType == "Pedido concluído" {
		return echo.NewHTTPError(http.StatusInternalServerError, "Sales order cant be updated")
	}

	var id_check int = 0
	query := `SELECT id FROM sales_orders_itens WHERE id_accessory = $1 AND id_sales_order = $2`
	db.QueryRow(query, id_accessory, id).Scan(&id_check)

	if id_check == 0 {
		query := `INSERT INTO sales_orders_itens(id_accessory, id_sales_order, qty, unit_price) VALUES ($1, $2, $3, $4)`

		_, err := db.Exec(query, id_accessory, id, 1, acc.PriceSell)
		if err != nil {
			return err
		}
	} else {
		return echo.NewHTTPError(http.StatusInternalServerError, "Accessory already added to this sales order")
	}

	// Don't check log insert error for now
	_ = audit.InsertLog(id_user, "/orders/:id/accessory/:id_accessory", "UPDATE", fmt.Sprintf("UPDATE request on sales order of id %d, add accessory", id), query)

	return nil
}

func UpdateNegotiationAdvanceStage(id int, id_user int, negT *models.AdvanceNegotiationRequest) error {
	db := storage.GetDB()

	query := `UPDATE so_business SET `
	params := []interface{}{}
	paramCount := 0

	if *negT.NewStage > 5 || *negT.NewStage < 1 {
		return echo.NewHTTPError(http.StatusInternalServerError, "New stage value not allowed")
	}

	if negT.NewStage != nil && negT.CurrentStage != nil {
		paramCount++
		query += fmt.Sprintf("stage = $%d, stage_last_updated_at = CURRENT_TIMESTAMP", paramCount)
		params = append(params, *&negT.NewStage)
	}

	if len(params) == 0 {
		return nil
	}

	paramCount++
	query += fmt.Sprintf(" WHERE id = $%d", paramCount)
	params = append(params, id)

	_, err := db.Exec(query, params...)
	if err != nil {
		return err
	}

	var id_customer *int = nil
	query = `SELECT id_customer FROM so_business WHERE id = $1`
	db.QueryRow(query, id).Scan(&id_customer)

	query = "INSERT INTO business_histories (id_user, id_customer, description, stage, id_mean_communication, id_business) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err = db.Exec(query, id_user, id_customer, "AVANÇO DE FUNIL", negT.NewStage, 4, id)
	if err != nil {
		return err
	}

	// Don't check log insert error for now
	_ = audit.InsertLog(id_user, "/negotiations/:id/advance", "UPDATE", fmt.Sprintf("UPDATE request on business %d, advance stage", id), query)

	return nil
}
