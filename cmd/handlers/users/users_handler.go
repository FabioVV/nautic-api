package users

import (
	"nautic/cmd/repositories"
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

	if err := repositories.InsertUser(user); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "user created successfully",
	})
}

func GetUser(c echo.Context) error {
	idParam := c.Param("id")

	userID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID format")
	}

	user, err := repositories.GetUser(userID)
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

	users, numRecords, err := repositories.GetUsers(qpage, qperpage, qname, qemail, qactive)
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

	if err := c.Validate(user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": validation.FmtErrReturn(err)})
	}

	userID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID format")
	}

	err = repositories.UpdateUser(userID, user)
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

	err = repositories.DeactivateUser(userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "user deactivated successfully",
	})
}
