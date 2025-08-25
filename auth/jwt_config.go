package auth

import (
	"os"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)



type JwtCustomClaims struct {
	Name string `json:"name"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}


func GetJwtConfig() echojwt.Config {
	var config = echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims{
			return new(JwtCustomClaims)
		},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}


	return config;
}

