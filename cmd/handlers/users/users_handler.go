package users

import (
	"nautic/cmd/storage"
	"nautic/cmd/utils"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func InsertUser(c echo.Context) error {
	db := storage.GetDB()

	name := c.FormValue("name")
	email := c.FormValue("email")
	phone := c.FormValue("phone")
	password := c.FormValue("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := "INSERT INTO users (name, email, phone, password_hash) VALUES ($1, $2, $3, $4)"

	_, err = db.Exec(query, name, email, phone, hashedPassword)
	if err != nil {
		if errU, ok := utils.CheckForUserError("email_unique", err); ok {
			return echo.NewHTTPError(errU.HttpErrCode, errU)
		}
		return err
	}

	return nil
}
