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
	"nautic/cmd/handlers/products"
	"nautic/cmd/handlers/users"
	"nautic/cmd/handlers/sales"

	nmiddleware "nautic/cmd/middleware"
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

	apiv1 := e.Group("/api/v1")

	/*AUTHENTICATION ROUTES*/
	authRoutes := apiv1.Group("/auth")
	authRoutes.POST("/signin", auth_h.Login)
	/*AUTHENTICATION ROUTES*/

	/*USER ROUTES*/
	userRoutes := apiv1.Group("/users")
	userRoutes.Use(echojwt.WithConfig(configJwt))
	userRoutes.Use(nmiddleware.CheckRoleAndPermissions)

	userRoutes.POST("", users.InsertUser)
	userRoutes.GET("", users.GetUsers)
	userRoutes.GET("/:id", users.GetUser)
	userRoutes.PATCH("/:id", users.UpdateUser)
	userRoutes.DELETE("/:id", users.DeactivateUser)
	/*USER ROUTES*/

	/*PERMISSIONS/ROLES ROUTES*/
	permsRoutes := apiv1.Group("/permissions")
	permsRoutes.Use(echojwt.WithConfig(configJwt))
	permsRoutes.Use(nmiddleware.CheckRoleAndPermissions)
	permsRoutes.GET("", auth_h.GetPermissions)

	/*PERMISSIONS/ROLES ROUTES*/

	/*ACESSORIES ROUTES*/
	accRoutes := apiv1.Group("/accessories")
	accRoutes.Use(echojwt.WithConfig(configJwt))
	accRoutes.Use(nmiddleware.CheckRoleAndPermissions)

	accRoutes.GET("", products.GetAccessories)
	accRoutes.POST("", products.InsertAccessory)
	accRoutes.DELETE("/:id", products.DeactivateAccessory)
	accRoutes.GET("/:id", products.GetAccessory)
	accRoutes.PATCH("/:id", products.UpdateAccessory)

	//I know, i know.... this should problably a separate group, but you know how things are.....
	accRoutes.GET("/accessories-types", products.GetAccessoriesTypes)
	accRoutes.POST("/accessories-types", products.InsertAccessoryType)
	accRoutes.DELETE("/accessories-types/:id", products.DeactivateAccessoryType)
	accRoutes.PATCH("/accessories-types/:id", products.UpdateAccessoryType)
	/*ACESSORIES ROUTES*/

	/*SALES ROUTES*/
	salRoutes := apiv1.Group("/sales")
	salRoutes.Use(echojwt.WithConfig(configJwt))
	salRoutes.Use(nmiddleware.CheckRoleAndPermissions)

	salRoutes.GET("/communication-means", sales.GetComMeans)
	salRoutes.POST("/communication-means", sales.InsertComMeans)
	salRoutes.DELETE("/communication-means/:id", sales.DeactivateComMeans)
	salRoutes.PATCH("/communication-means/:id", sales.UpdateComMeans)

	/*SALES ROUTES*/



	e.Logger.Fatal(e.Start(":8080"))
}
