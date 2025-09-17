package products

import (
	"nautic/cmd/repositories"
	"nautic/models"
	"nautic/validation"
	"strconv"

	"net/http"

	"github.com/labstack/echo/v4"
)

func UpdateAccessory(c echo.Context) error {
	idParam := c.Param("id")

	accT := new(models.UpdateAccessoryRequest)

	if err := c.Bind(accT); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := c.Validate(accT); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": validation.FmtErrReturn(err)})
	}

	accTID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	err = repositories.UpdateAccessory(accTID, accT)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Accessory updated successfully",
	})
}

func UpdateAccessoryType(c echo.Context) error {
	idParam := c.Param("id")

	accT := new(models.UpdateAccessoryTypeRequest)

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

	err = repositories.UpdateAccessoryType(accTID, accT)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Accessory type updated successfully",
	})
}

func DeactivateAccessoryType(c echo.Context) error {
	idParam := c.Param("id")

	accTID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	err = repositories.DeactivateAccessoryType(accTID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusNoContent, echo.Map{
		"message": "type deactivated successfully",
	})
}

func InsertAccessoryType(c echo.Context) error {
	accT := new(models.CreateAccessoryTypeRequest)

	if err := c.Bind(accT); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := c.Validate(accT); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": validation.FmtErrReturn(err)})
	}

	if err := repositories.InsertAccessoryType(accT); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "accessory type created successfully",
	})
}

func DeactivateAccessory(c echo.Context) error {
	idParam := c.Param("id")

	accID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	err = repositories.DeactivateAccessory(accID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusNoContent, echo.Map{
		"message": "accessory deactivated successfully",
	})
}

func InsertAccessory(c echo.Context) error {
	accT := new(models.CreateAccessoryRequest)

	if err := c.Bind(accT); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload"+err.Error())
	}

	if err := c.Validate(accT); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": validation.FmtErrReturn(err)})
	}

	if err := repositories.InsertAccessory(accT); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "accessory type created successfully",
	})
}

func GetAccessories(c echo.Context) error {
	qpage := c.QueryParams().Get("pageNumber")
	qperpage := c.QueryParams().Get("perPage")
	qmodel := c.QueryParams().Get("name")
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

func GetAccessory(c echo.Context) error {
	idParam := c.Param("id")

	accID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	acc, err := repositories.GetAccessory(accID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": acc,
	})
}

func GetAccessoriesTypes(c echo.Context) error {

	qpage := c.QueryParams().Get("pageNumber")
	qperpage := c.QueryParams().Get("perPage")
	qtype := c.QueryParams().Get("type")
	qactive := c.QueryParams().Get("active")

	accsT, numRecords, err := repositories.GetAccessoriesTypes(qpage, qperpage, qtype, qactive)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":         accsT,
		"totalRecords": numRecords,
	})
}
