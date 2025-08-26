package repositories

import (
	"database/sql"
	"nautic/cmd/storage"
	"nautic/cmd/utils"
	"nautic/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func GetUser(id int) (models.User, error) {
	db := storage.GetDB()

	var user models.User
	query := `SELECT id, name, email, active, phone, created_at, updated_at FROM users WHERE id = $1`

	err := db.QueryRow(query, id).Scan(&user.Id, &user.Name, &user.Email, &user.Active, &user.Phone, &user.CreatedAt, &user.UpdatedAt)

	if err := db.QueryRow(query, id).Scan(&user.Id, &user.Name, &user.Email, &user.Active, &user.Phone, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return user, echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return user, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve user")
	}

	return user, err
}

func InsertUser(user *models.CreateUserRequest) error {
	db := storage.GetDB()

	if errMsg, ok := utils.IsGoodPassword(user.Password); !ok {
		return echo.NewHTTPError(http.StatusBadRequest, errMsg)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := "INSERT INTO users (name, email, phone, password_hash) VALUES ($1, $2, $3, $4)"

	_, err = db.Exec(query, user.Name, user.Email, user.Phone, hashedPassword)
	if err != nil {
		if errU, ok := utils.CheckForUserError("unique_email", err); ok {
			return echo.NewHTTPError(errU.HttpErrCode, errU)
		}
		return err
	}

	return nil
}

func UpdateUser(id int, user *models.UpdateUserRequest) error {
	//db := storage.GetDB()

	return nil

}
