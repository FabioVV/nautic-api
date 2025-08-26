package users

import (
	"nautic/cmd/storage"
	"nautic/cmd/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func InsertUser(c echo.Context) error {
	db := storage.GetDB()


	name := c.FormValue("name")
	email := c.FormValue("email")
	phone := c.FormValue("phone")
	password := c.FormValue("password")

	if !utils.IsGoodText(name) {
		return echo.NewHTTPError(http.StatusBadRequest, "Name is required")
	}

	if !utils.IsValidEmail(email) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid email address")
	}

	if errMsg, ok := utils.IsGoodPassword(password); !ok {
		return echo.NewHTTPError(http.StatusBadRequest, errMsg)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := "INSERT INTO users (name, email, phone, password_hash) VALUES ($1, $2, $3, $4)"

	_, err = db.Exec(query, name, email, phone, hashedPassword)
	if err != nil {
		if errU, ok := utils.CheckForUserError("unique_email", err); ok {
			return echo.NewHTTPError(errU.HttpErrCode, errU)
		}
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "user created successfully",
		"status": "success",
	})
}
