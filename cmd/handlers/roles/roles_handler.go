package roles

import (
	"nautic/cmd/repositories"
	"nautic/cmd/utils"
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

	claims, err := utils.GetLoggedInUserClaims(c)
	if err != nil {
		return err
	}

	if err := repositories.InsertRole(role, claims.Id); err != nil {
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
	qshow_admin := c.QueryParams().Get("show_admin")

	claims, err := utils.GetLoggedInUserClaims(c)
	if err != nil {
		return err
	}

	roles, numRecords, err := repositories.GetRoles(qpage, qperpage, qname, qshow_admin, claims.Id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":         roles,
		"totalRecords": numRecords,
	})
}

func GetRolePermissions(c echo.Context) error {
	idParam := c.Param("id")

	rID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	roles, err := repositories.GetRolePermissions(rID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": roles,
	})
}

func GetRole(c echo.Context) error {
	idParam := c.Param("id")

	rID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	claims, err := utils.GetLoggedInUserClaims(c)
	if err != nil {
		return err
	}

	role, err := repositories.GetRole(rID, claims.Id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, role)
}

func RemoveRolePermission(c echo.Context) error {
	idParam := c.Param("id")
	idParamPerm := c.Param("id_perm")

	roleId, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	permId, err := strconv.Atoi(idParamPerm)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	claims, err := utils.GetLoggedInUserClaims(c)
	if err != nil {
		return err
	}

	err = repositories.RemoveRolePermission(roleId, permId, claims.Id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "role permission removed successfully",
	})
}

func DeleteRole(c echo.Context) error {
	idParam := c.Param("id")

	rID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID format")
	}

	claims, err := utils.GetLoggedInUserClaims(c)
	if err != nil {
		return err
	}

	err = repositories.DeleteRole(rID, claims.Id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "role deleted successfully",
	})
}

func InsertRolePermission(c echo.Context) error {
	idParam := c.Param("id")
	idParamPerm := c.Param("id_perm")

	roleId, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	permId, err := strconv.Atoi(idParamPerm)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	claims, err := utils.GetLoggedInUserClaims(c)
	if err != nil {
		return err
	}

	err = repositories.InsertRolePermission(roleId, permId, claims.Id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "role updated successfully",
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

	claims, err := utils.GetLoggedInUserClaims(c)
	if err != nil {
		return err
	}

	err = repositories.UpdateRole(rID, claims.Id, rT)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "role updated successfully",
	})
}
