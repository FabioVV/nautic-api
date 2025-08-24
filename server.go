package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"

	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"nautic/cmd/storage"
	nauticmiddleware "nautic/cmd/middleware"
	nauticjwt "nautic/auth"

)



func secretRoute(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*nauticjwt.JwtCustomClaims)
	name := claims.Name
	admin := claims.Admin

	if admin{
		return c.String(http.StatusOK, "Welcome, "+name+"!\n"+"You are an admin")
	}

	return c.String(http.StatusOK, "Welcome, "+name+"!\n"+"You are not an admin")
}

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username != "fabio" || password != "shhh!"{
		return echo.ErrUnauthorized
	}

	claims := &nauticjwt.JwtCustomClaims{
		"fabio",
		true,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil{
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})

}

func main(){
	e := echo.New()

	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
		return
	}

	storage.InitDB()

	//e.Use(middleware.Logger())
	e.Use(nauticmiddleware.LogReq)
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"localhost:8080", "localhost:8081"},
		AllowHeaders: []string{},
	}))

	authRoutes := e.Group("/auth")
	authRoutes.POST("/login", login)


	r := e.Group("/restricted")
	r.Use(echojwt.WithConfig(nauticjwt.GetJwtConfig()))

	r.GET("/secret", secretRoute)

	e.Logger.Fatal(e.Start(":8080"))
}
