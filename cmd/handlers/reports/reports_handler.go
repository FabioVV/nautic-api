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

	qdateIni := c.QueryParams().Get("dateIni")
	qdateEnd := c.QueryParams().Get("dateEnd")

	data, numRecords, err := repositories.GetNegotiationsReport(qpage, qperpage, qname, qpboat, qdateIni, qdateEnd)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":         data,
		"totalRecords": numRecords,
	})
}

func GetSalesOrdersReport(c echo.Context) error {

	qpage := c.QueryParams().Get("pageNumber")
	qperpage := c.QueryParams().Get("perPage")
	qname := c.QueryParams().Get("name")
	qseller := c.QueryParams().Get("seller")

	qdateIni := c.QueryParams().Get("dateIni")
	qdateEnd := c.QueryParams().Get("dateEnd")

	data, numRecords, err := repositories.GetSalesOrdersReport(qpage, qperpage, qname, qseller, qdateIni, qdateEnd)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":         data,
		"totalRecords": numRecords,
	})
}

func GetLostNegotiationsReport(c echo.Context) error {

	qpage := c.QueryParams().Get("pageNumber")
	qperpage := c.QueryParams().Get("perPage")
	qname := c.QueryParams().Get("name")
	qpboat := c.QueryParams().Get("boat")

	qdateIni := c.QueryParams().Get("dateIni")
	qdateEnd := c.QueryParams().Get("dateEnd")

	data, numRecords, err := repositories.GetLostNegotiationsReport(qpage, qperpage, qname, qpboat, qdateIni, qdateEnd)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{
		"data":         data,
		"totalRecords": numRecords,
	})
}
