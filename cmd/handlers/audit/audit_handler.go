package audit

import (
	"nautic/cmd/repositories"
	"nautic/cmd/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetLogs(c echo.Context) error {

	qpage := c.QueryParams().Get("pageNumber")
	qperpage := c.QueryParams().Get("perPage")
	qaction := c.QueryParams().Get("action")
	qdesc := c.QueryParams().Get("description")

	claims, err := utils.GetLoggedInUserClaims(c)
	if err != nil {
		return err
	}

	logs, numRecords, err := repositories.GetLogs(claims.Id, qpage, qperpage, qaction, qdesc)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":         logs,
		"totalRecords": numRecords,
	})
}
