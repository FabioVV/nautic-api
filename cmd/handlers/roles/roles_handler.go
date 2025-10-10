package roles

import (
	"nautic/cmd/repositories"
	"nautic/models"
	"nautic/validation"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func InsertRole(c echo.Context) error {
	role := new(models.CreateRole)

	if err := c.Bind(role); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := c.Validate(role); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": validation.FmtErrReturn(err)})
	}

	if err := repositories.InsertRole(role); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "role created successfully",
	})
}

func GetRoles(c echo.Context) error {

	qpage := c.QueryParams().Get("pageNumber")
	qperpage := c.QueryParams().Get("perPage")
	qname := c.QueryParams().Get("name")

	roles, numRecords, err := repositories.GetRoles(qpage, qperpage, qname)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":         roles,
		"totalRecords": numRecords,
	})
}

func GetRole(c echo.Context) error {
	idParam := c.Param("id")

	rID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID format")
	}

	role, err := repositories.GetRole(rID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, role)
}

func DeleteRole(c echo.Context) error {
	idParam := c.Param("id")

	rID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID format")
	}

	err = repositories.DeleteRole(rID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "role deleted successfully",
	})
}

func UpdateRole(c echo.Context) error {
	idParam := c.Param("id")

	rT := new(models.CreateRole)

	if err := c.Bind(rT); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := c.Validate(rT); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": validation.FmtErrReturn(err)})
	}

	rID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID format")
	}

	err = repositories.UpdateRole(rID, rT)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "role updated successfully",
	})
}
