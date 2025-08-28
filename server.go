package main

import (
	"log"

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
	vali.RegisterTagNameFunc(validation.GetJsonStructName)
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
		AllowOrigins: []string{"localhost:8080", "http://localhost:4200", "localhost:4200", "http://127.0.0.1:4200"},
	}))

	e.POST("/tmpr", users.InsertUser)

	apiv1 := e.Group("/api/v1")

	authRoutes := apiv1.Group("/auth")
	authRoutes.POST("/signin", auth_h.Login)

	userRoutes := apiv1.Group("/users")
	userRoutes.Use(echojwt.WithConfig(configJwt))

	userRoutes.POST("", users.InsertUser)
	userRoutes.GET("/:id", users.GetUser)
	userRoutes.PATCH("/:id", users.UpdateUser)
	userRoutes.DELETE("/:id", users.DeactivateUser)

	e.Logger.Fatal(e.Start(":8080"))
}
