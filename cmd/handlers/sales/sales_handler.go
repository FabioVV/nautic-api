package sales

import (
	"nautic/cmd/repositories"
	"nautic/cmd/utils"
	"nautic/models"
	"nautic/validation"
	"strconv"

	"net/http"

	"github.com/labstack/echo/v4"
)

func GetNegotiation(c echo.Context) error {
	idParam := c.Param("id")

	negID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	acc, err := repositories.GetNegotiation(negID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": acc,
	})
}

func GetNegotiationHistory(c echo.Context) error {
	idParam := c.Param("id")

	claims, err := utils.GetLoggedInUserClaims(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	data, numRecords, err := repositories.GetNegotiationHistory(id, claims.Id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":         data,
		"totalRecords": numRecords,
	})
}

func GetNegotiations(c echo.Context) error {

	claims, err := utils.GetLoggedInUserClaims(c)
	if err != nil {
		return err
	}

	// qpage := c.QueryParams().Get("pageNumber")
	// qperpage := c.QueryParams().Get("perPage")
	qsearch := c.QueryParams().Get("search")
	//qactive := c.QueryParams().Get("active")

	data, numRecords, err := repositories.GetNegotiations(qsearch, claims.Id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":         data,
		"totalRecords": numRecords,
	})
}

func GetCustomersBirthday(c echo.Context) error {
	data, numRecords, err := repositories.GetCustomersBirthday()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":         data,
		"totalRecords": numRecords,
	})
}

func GetCustomers(c echo.Context) error {

	qpage := c.QueryParams().Get("pageNumber")
	qperpage := c.QueryParams().Get("perPage")
	qname := c.QueryParams().Get("name")
	qemail := c.QueryParams().Get("email")
	qphone := c.QueryParams().Get("phone")
	qpboat := c.QueryParams().Get("boat")

	data, numRecords, err := repositories.GetCustomers(qpage, qperpage, qname, qemail, qphone, qpboat)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":         data,
		"totalRecords": numRecords,
	})
}

func GetComMeans(c echo.Context) error {

	qpage := c.QueryParams().Get("pageNumber")
	qperpage := c.QueryParams().Get("perPage")
	qname := c.QueryParams().Get("name")
	qactive := c.QueryParams().Get("active")

	data, numRecords, err := repositories.GetComMeans(qpage, qperpage, qname, qactive)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":         data,
		"totalRecords": numRecords,
	})
}

func UpdateComMeans(c echo.Context) error {
	idParam := c.Param("id")

	accT := new(models.UpdateCommunicationMeaneRequest)

	if err := c.Bind(accT); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := c.Validate(accT); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": validation.FmtErrReturn(err)})
	}

	accTID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID format")
	}

	err = repositories.UpdateComMean(accTID, accT)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "communication mean updated successfully",
	})
}

func DeactivateComMeans(c echo.Context) error {
	idParam := c.Param("id")

	accTID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	err = repositories.DeactivateComMean(accTID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "communication mean deactivated successfully",
	})
}

func InsertComMeans(c echo.Context) error {
	accT := new(models.CreateCommunicationMeanRequest)

	if err := c.Bind(accT); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := c.Validate(accT); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": validation.FmtErrReturn(err)})
	}

	if err := repositories.InsertComMeans(accT); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "communication mean created successfully",
	})
}

func InsertNegotiationHistory(c echo.Context) error {
	idParam := c.Param("id")

	claims, err := utils.GetLoggedInUserClaims(c)
	if err != nil {
		return err
	}

	negT := new(models.CreateNegotiationHistoryRequest)

	if err := c.Bind(negT); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := c.Validate(negT); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": validation.FmtErrReturn(err)})
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	if claims.Id != int(*negT.UserId) {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid ID for this resource")
	}

	err = repositories.CreateNegotiationHistory(id, negT)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Negotiation history created successfully",
	})
}

func InsertNegotiation(c echo.Context) error {
	neg := new(models.CreateNegotiationRequest)

	if err := c.Bind(neg); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := c.Validate(neg); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": validation.FmtErrReturn(err)})
	}

	if err := repositories.InsertNegotiation(neg); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Negotiation created successfully",
	})
}

func UpdateNegotiation(c echo.Context) error {
	idParam := c.Param("id")

	negT := new(models.CreateNegotiationRequest)

	if err := c.Bind(negT); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload"+err.Error())
	}

	if err := c.Validate(negT); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": validation.FmtErrReturn(err)})
	}

	accTID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	err = repositories.UpdateNegotiation(accTID, negT)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Engine updated successfully",
	})
}
