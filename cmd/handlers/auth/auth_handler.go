package auth_h

import (
	"database/sql"
	"nautic/auth"
	r "nautic/cmd/repositories"
	"nautic/cmd/storage"
	"nautic/cmd/utils"
	"nautic/models"
	"nautic/validation"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var RoutesPermissions = map[string]string{
	"POST:/api/v1/users":       "users:create",
	"GET:/api/v1/users":        "users:view",
	"GET:/api/v1/users/:id":    "users:view",
	"PATCH:/api/v1/users/:id":  "users:update",
	"DELETE:/api/v1/users/:id": "users:delete",
}

func CheckRoleAndPermissions(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse user credentials"})
		}
		claims, ok := user.Claims.(*auth.JwtCustomClaims)
		if !ok {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse user credentials claims"})
		}

		if len(claims.Roles) > 0 && claims.Roles[0] == "admin" {
			return next(c)
		}

		routePermissionKey := c.Request().Method + ":" + c.Path()
		routePermission := RoutesPermissions[routePermissionKey]

		for _, perm := range claims.Permissions {
			if routePermission == perm {
				next(c)
			}
		}

		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "User does not have permission for the requested resource"})
	}
}

func Login(c echo.Context) error {
	db := storage.GetDB()

	lr := new(models.LoginRequest)
	var name string
	var email string
	var password string
	var id int

	if err := c.Bind(lr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := c.Validate(lr); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": validation.FmtErrReturn(err)})
	}

	query := `SELECT id, name, email, password_hash FROM users WHERE email = $1`
	err := db.QueryRow(query, lr.Email).Scan(&id, &name, &email, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error (db)")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(lr.Password)); err != nil {
		if errU, ok := utils.GetUserError(err.Error()); ok {
			return echo.NewHTTPError(errU.HttpErrCode, errU)
		}
		return err
	}

	userRoles, err := r.GetUserRoles(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve user roles during auth: "+err.Error())
	}

	userPermissions, err := r.GetUserPermissions(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve user permissions during auth: "+err.Error())
	}

	claims := &auth.JwtCustomClaims{
		name,
		userRoles,
		userPermissions,
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
		"name":  name,
		"email": email,
	})

}
