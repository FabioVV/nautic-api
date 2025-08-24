package auth

import (
	"os"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)


type JwtCustomClaims struct {
	Name string `json:"name"`
	Admin bool `json:"admin"`
	jwt.RegisteredClaims
}

var config = echojwt.Config{
	NewClaimsFunc: func(c echo.Context) jwt.Claims{
		return new(JwtCustomClaims)
	},
	SigningKey: []byte(os.Getenv("JWT_SECRET")),
}

func GetJwtConfig() echojwt.Config {
	return config;
}

