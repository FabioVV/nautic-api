package main

import (
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"nautic/auth"
	auth_h "nautic/cmd/handlers/auth"
	"nautic/cmd/handlers/users"
	"nautic/cmd/storage"
	"nautic/validation"
)

func main() {
	e := echo.New()
	vali := validator.New()
	vali.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	e.Validator = &validation.CustomValidator{Validator: vali}

	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}

	storage.InitDB()
	defer storage.CloseDB()

	configJwt := auth.GetJwtConfig()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"localhost:8080", "localhost:8081"},
		AllowHeaders: []string{},
	}))

	e.POST("/tmpr", users.InsertUser)

	authRoutes := e.Group("/auth")
	authRoutes.POST("/signin", auth_h.Login)

	userRoutes := e.Group("/users")
	userRoutes.Use(echojwt.WithConfig(configJwt))

	userRoutes.POST("", users.InsertUser)
	userRoutes.GET("/:id", users.GetUser)
	userRoutes.PATCH("/:id", users.UpdateUser)
	userRoutes.DELETE("/:id", users.DeactivateUser)

	e.Logger.Fatal(e.Start(":8080"))
}
