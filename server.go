package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"nautic/auth"
	"nautic/cmd/handlers/auth"
	"nautic/cmd/handlers/users"
	"nautic/cmd/storage"
)

func main() {
	e := echo.New()

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

	e.Logger.Fatal(e.Start(":8080"))
}
