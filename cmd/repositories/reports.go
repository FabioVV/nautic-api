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

func GetNegotiationsReport(pagenum string, limitPerPage string, name string, boat string) ([]models.NegotiationReport, int, error) {
	db := storage.GetDB()

	var negs []models.NegotiationReport

	pagenumber, err := strconv.Atoi(pagenum)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve report (PG1)")
	}
	limit, err := strconv.Atoi(limitPerPage)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve report (PG2)")
	}

	offset := (pagenumber - 1) * limit

	conds := []string{}
	args := []interface{}{}
	paramCount := 1

	if name != "" {
		conds = append(conds, fmt.Sprintf("C.name ILIKE $%d", paramCount))
		args = append(args, "%"+name+"%")
		paramCount++
	}

	if boat != "" {
		conds = append(conds, fmt.Sprintf("SB.boat_name ILIKE $%d", paramCount))
		args = append(args, "%"+boat+"%")
		paramCount++
	}

	conds = append(conds, fmt.Sprintf("SB.negotiation_active = $%d", paramCount))
	args = append(args, "Y")
	paramCount++

	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}

	args = append(args, limitPerPage, offset)
	limitArgPos := paramCount
	offsetArgPos := paramCount + 1

	query := fmt.Sprintf(`
		SELECT
		SB.id,
		SB.id_customer,
		SB.id_mean_communication,
		MC.name AS mean_communication_name,
		C.name,
		C.email,
		C.phone,
		SB.boat_name,
		SB.estimated_value,
		SB.stage,
		
		-- days since stage_last_updated_at (integer days)
		DATE_PART('day', NOW() - SB.stage_last_updated_at)::bigint AS days_since_stage_change,

		-- most recent history timestamp per business
		BH.last_history_at,

		-- days since last history (integer days), null if no history
		CASE WHEN BH.last_history_at IS NOT NULL
			THEN DATE_PART('day', NOW() - BH.last_history_at)::bigint
			ELSE NULL END AS days_since_last_history

		-- full interval since last history (null if no history)
		-- CASE WHEN BH.last_history_at IS NOT NULL
		--	THEN NOW() - BH.last_history_at
		--	ELSE NULL END AS time_since_last_history

		FROM so_business AS SB
		INNER JOIN customers AS C ON SB.id_customer = C.id
		LEFT JOIN mean_communication AS MC ON SB.id_mean_communication = MC.id
		-- subquery to get most recent history per business (id_business links to so_business.id)
		LEFT JOIN (
		SELECT id_business, MAX(created_at) AS last_history_at
		FROM business_histories
		GROUP BY id_business
		) AS BH ON BH.id_business = SB.id
	%s
	ORDER BY SB.id
	LIMIT $%d OFFSET $%d
	`, where, limitArgPos, offsetArgPos)

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
		LEFT JOIN customers AS C ON SB.id_customer = C.id
		LEFT JOIN mean_communication AS MC ON SB.id_mean_communication = MC.id
		-- subquery to get most recent history per business (id_business links to so_business.id)
		LEFT JOIN (
		SELECT id_business, MAX(created_at) AS last_history_at
		FROM business_histories
		GROUP BY id_business
		) AS BH ON BH.id_business = SB.id
	%s
	`, where)
	//println(queryTotalRecords)

	rowsCount := db.QueryRow(queryTotalRecords, args...)
	numRecords := 0
	rowsCount.Scan(&numRecords)

	for rows.Next() {
		var curNeg models.NegotiationReport

		if err := rows.Scan(&curNeg.Id, &curNeg.CustomerId, &curNeg.MeanComId,
			&curNeg.Name, &curNeg.Email, &curNeg.Phone, &curNeg.MeamComName,
			&curNeg.BoatName, &curNeg.EstimatedValue, &curNeg.Stage, &curNeg.DaysSinceStageChange, &curNeg.LastHistoryAt, &curNeg.DaysSinceLastHistory); err != nil {
			return nil, 0, fmt.Errorf("scan error: %w", err)
		}

		negs = append(negs, curNeg)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return negs, numRecords, nil
}
