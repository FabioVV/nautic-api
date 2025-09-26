package utils

import (
	"nautic/auth"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GetLoggedInUserClaims(c echo.Context) (*auth.JwtCustomClaims, error) {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return nil, c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse user credentials"})
	}
	claims, ok := user.Claims.(*auth.JwtCustomClaims)
	if !ok {
		return nil, c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse user credentials claims"})
	}

	return claims, nil
}
