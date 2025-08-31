package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"nautic/auth"
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
