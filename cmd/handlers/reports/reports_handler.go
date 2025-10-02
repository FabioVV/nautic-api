package reports

import (
	"nautic/cmd/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetNegotiationsReport(c echo.Context) error {

	qpage := c.QueryParams().Get("pageNumber")
	qperpage := c.QueryParams().Get("perPage")
	qname := c.QueryParams().Get("name")
	qpboat := c.QueryParams().Get("boat")

	data, numRecords, err := repositories.GetNegotiationsReport(qpage, qperpage, qname, qpboat)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":         data,
		"totalRecords": numRecords,
	})
}
