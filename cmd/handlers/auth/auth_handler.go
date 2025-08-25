package auth_h

import (
	"database/sql"
	"nautic/auth"
	"nautic/cmd/storage"
	"nautic/models"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Login(c echo.Context) error {
	db := storage.GetDB()
	if db == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Connection failed")
	}

	var user models.User
	email := c.FormValue("email")
	password := c.FormValue("password")

	query := `SELECT name, password_hash FROM users WHERE email = $1`
	err := db.QueryRow(query, email).Scan(&user.Name, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error (db)")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return err
	}

	claims := &auth.JwtCustomClaims{
		user.Name,
		"No role yet",
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})

}
