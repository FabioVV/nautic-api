package main

import (
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"nautic/auth"
	auth_h "nautic/cmd/handlers/auth"
	"nautic/cmd/handlers/users"
	"nautic/cmd/storage"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

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
	userRoutes.PATCH("/:id", users.InsertUser)
	userRoutes.DELETE("/:id", users.InsertUser)

	e.Logger.Fatal(e.Start(":8080"))
}
