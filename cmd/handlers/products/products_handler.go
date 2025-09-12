package products

import (
	"nautic/cmd/repositories"

	"net/http"

	"github.com/labstack/echo/v4"
)

func GetAccessories(c echo.Context) error {

	qpage := c.QueryParams().Get("pageNumber")
	qperpage := c.QueryParams().Get("perPage")
	qmodel := c.QueryParams().Get("model")
	qactive := c.QueryParams().Get("active")

	accs, numRecords, err := repositories.GetAccessories(qpage, qperpage, qmodel, qactive)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":         accs,
		"totalRecords": numRecords,
	})
}
