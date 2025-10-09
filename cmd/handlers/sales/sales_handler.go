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

func GetCustomer(c echo.Context) error {
	idParam := c.Param("id")

	cID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	acc, err := repositories.GetCustomer(cID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": acc,
	})
}

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

func GetNegotiationsAlerts(c echo.Context) error {

	claims, err := utils.GetLoggedInUserClaims(c)
	if err != nil {
		return err
	}

	data, err := repositories.GetNegotiationsAlerts(claims.Id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": data,
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

func GetCustomerNegotiationsHistories(c echo.Context) error {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	data, numRecords, err := repositories.GetCustomerNegotiationsHistories(id)
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

func UpdateBoatSalesOrder(c echo.Context) error {
	idParam := c.Param("id")
	idParamBoat := c.Param("id_boat")

	soID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	boatID, err := strconv.Atoi(idParamBoat)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	err = repositories.UpdateBoatSalesOrder(soID, boatID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Sales order updated successfully",
	})
}

func UpgradeQuoteToSalesOrder(c echo.Context) error {
	idParam := c.Param("id")

	soID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	err = repositories.UpgradeQuoteToSalesOrder(soID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Sales order updated successfully",
	})
}

func CancelSalesOrder(c echo.Context) error {
	idParam := c.Param("id")

	soID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	err = repositories.CancelSalesOrder(soID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Sales order updated successfully",
	})
}

func UpdateEngineSalesOrder(c echo.Context) error {
	idParam := c.Param("id")
	idParamEngine := c.Param("id_engine")

	soID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	engID, err := strconv.Atoi(idParamEngine)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	err = repositories.UpdateEngineSalesOrder(soID, engID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Sales order updated successfully",
	})
}

func RemoveAccessorySalesOrder(c echo.Context) error {
	idParam := c.Param("id")
	idParamAccessory := c.Param("id_accessory")

	soID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	accID, err := strconv.Atoi(idParamAccessory)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	err = repositories.RemoveAccessorySalesOrder(soID, accID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Sales order updated successfully",
	})
}

func UpdateAccessoryQtySalesOrder(c echo.Context) error {
	idParam := c.Param("id")
	idParamAccessory := c.Param("id_accessory")

	soID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	accID, err := strconv.Atoi(idParamAccessory)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	accQty := new(models.UpdateSalesOrderItemQty)

	if err := c.Bind(accQty); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := c.Validate(accQty); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": validation.FmtErrReturn(err)})
	}

	err = repositories.UpdateAccessoryQtySalesOrder(soID, accID, accQty)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Sales order updated successfully",
	})
}

func UpdateAccessorySalesOrder(c echo.Context) error {
	idParam := c.Param("id")
	idParamAccessory := c.Param("id_accessory")

	soID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	accID, err := strconv.Atoi(idParamAccessory)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	err = repositories.UpdateAccessorySalesOrder(soID, accID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Sales order updated successfully",
	})
}

func InsertSalesOrder(c echo.Context) error {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	if err := repositories.InsertSalesOrderUsingBusinessHistory(id); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "sales order created successfully",
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

func UpdateCustomer(c echo.Context) error {
	idParam := c.Param("id")

	cT := new(models.CustomerRequest)

	if err := c.Bind(cT); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload"+err.Error())
	}

	if err := c.Validate(cT); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": validation.FmtErrReturn(err)})
	}

	accTID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	err = repositories.UpdateCustomer(accTID, cT)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Customer updated successfully",
	})
}

func ReactivateNegotiation(c echo.Context) error {
	idParam := c.Param("id")

	negID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	err = repositories.ReactivateNegotiation(negID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Negotiations updated successfully",
	})
}

func LostNegotiation(c echo.Context) error {
	idParam := c.Param("id")

	negT := new(models.LostNegotiation)

	if err := c.Bind(negT); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload"+err.Error())
	}

	if err := c.Validate(negT); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": validation.FmtErrReturn(err)})
	}

	negID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	err = repositories.LostNegotiation(negID, negT)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Negotiation updated successfully",
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
		"message": "Negotiation updated successfully",
	})
}

func UpdateNegotiationAdvanceStage(c echo.Context) error {
	idParam := c.Param("id")

	negT := new(models.AdvanceNegotiationRequest)

	claims, err := utils.GetLoggedInUserClaims(c)
	if err != nil {
		return err
	}

	if err := c.Bind(negT); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := c.Validate(negT); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": validation.FmtErrReturn(err)})
	}

	negID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	err = repositories.UpdateNegotiationAdvanceStage(negID, claims.Id, negT)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Negotiation updated successfully",
	})
}

func GetSalesOrder(c echo.Context) error {
	idParam := c.Param("id")

	idSalesOrder, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	data, err := repositories.GetSalesOrder(idSalesOrder)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": data,
	})
}

func GetSalesOrderItens(c echo.Context) error {
	idParam := c.Param("id")

	idSalesOrder, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	data, err := repositories.GetSalesOrderItens(idSalesOrder)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": data,
	})
}

func ChangeSalesOrderFileType(c echo.Context) error {
	idParam := c.Param("id")
	idFileParam := c.Param("id_file")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	idFile, err := strconv.Atoi(idFileParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file ID format")
	}

	soFT := new(models.UpdateSalesOrderFileType)

	if err := c.Bind(soFT); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := c.Validate(soFT); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": validation.FmtErrReturn(err)})
	}

	err = repositories.ChangeSalesOrderFileType(id, idFile, soFT)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Sales order file type changed successfully",
	})
}

func UploadSalesOrderFile(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}
	err = repositories.UploadSalesOrderFile(c, id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Sales order file uploaded successfully",
	})
}

func GetSalesOrderFiles(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	data, err := repositories.GetSalesOrderFiles(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": data,
	})
}

func RemoveSalesOrderFile(c echo.Context) error {
	idParam := c.Param("id")
	idFileParam := c.Param("id_file")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	idFile, err := strconv.Atoi(idFileParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file ID format")
	}

	err = repositories.RemoveSalesOrderFile(id, idFile)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Sales order file removed successfully",
	})
}
