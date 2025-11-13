package users

import (
	"nautic/cmd/repositories"
	"nautic/cmd/utils"

	"nautic/models"
	"nautic/validation"

	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func InsertUser(c echo.Context) error {
	user := new(models.CreateUserRequest)

	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := c.Validate(user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": validation.FmtErrReturn(err)})
	}

	claims, err := utils.GetLoggedInUserClaims(c)
	if err != nil {
		return err
	}

	if err := repositories.InsertUser(user, claims.Id); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "user created successfully",
	})
}

func GetUserPermissions(c echo.Context) error {
	idParam := c.Param("id")

	rID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	claims, err := utils.GetLoggedInUserClaims(c)
	if err != nil {
		return err
	}

	roles, err := repositories.GetUserPermissions__(rID, claims.Id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": roles,
	})
}

func RemoveUserPermission(c echo.Context) error {
	idParam := c.Param("id")
	idParamPerm := c.Param("id_perm")

	userId, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	permId, err := strconv.Atoi(idParamPerm)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	err = repositories.RemoveUserPermission(userId, permId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "user permission removed successfully",
	})
}

func InsertUserPermission(c echo.Context) error {
	idParam := c.Param("id")
	idParamPerm := c.Param("id_perm")

	userId, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	permId, err := strconv.Atoi(idParamPerm)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	err = repositories.InsertUserPermission(userId, permId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "user updated successfully",
	})
}

func GetUser(c echo.Context) error {
	idParam := c.Param("id")

	userID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID format")
	}

	claims, err := utils.GetLoggedInUserClaims(c)
	if err != nil {
		return err
	}

	user, err := repositories.GetUser(userID, claims.Id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func GetUsers(c echo.Context) error {

	qpage := c.QueryParams().Get("pageNumber")
	qperpage := c.QueryParams().Get("perPage")
	qname := c.QueryParams().Get("name")
	qemail := c.QueryParams().Get("email")
	qactive := c.QueryParams().Get("active")

	claims, err := utils.GetLoggedInUserClaims(c)
	if err != nil {
		return err
	}

	users, numRecords, err := repositories.GetUsers(claims.Id, qpage, qperpage, qname, qemail, qactive)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":         users,
		"totalRecords": numRecords,
	})
}

func UpdateUser(c echo.Context) error {
	idParam := c.Param("id")

	user := new(models.UpdateUserRequest)

	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	// if err := c.Validate(user); err != nil {
	// 	return c.JSON(http.StatusBadRequest, echo.Map{"errors": validation.FmtErrReturn(err)})
	// }

	userID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID format")
	}

	claims, err := utils.GetLoggedInUserClaims(c)
	if err != nil {
		return err
	}

	err = repositories.UpdateUser(userID, claims.Id, user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "user updated successfully",
	})
}

func DeactivateUser(c echo.Context) error {
	idParam := c.Param("id")

	userID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID format")
	}

	claims, err := utils.GetLoggedInUserClaims(c)
	if err != nil {
		return err
	}

	err = repositories.DeactivateUser(userID, claims.Id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "user deactivated successfully",
	})
}
