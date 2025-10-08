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
	"nautic/cmd/handlers/reports"
	"nautic/cmd/handlers/sales"
	"nautic/cmd/handlers/users"

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

	/*ENGINES ROUTES*/
	enginesRoutes := apiv1.Group("/engines")
	enginesRoutes.Use(echojwt.WithConfig(configJwt))
	enginesRoutes.Use(nmiddleware.CheckRoleAndPermissions)

	enginesRoutes.GET("", products.GetEngines)
	enginesRoutes.POST("", products.InsertEngine)
	enginesRoutes.DELETE("/:id", products.DeactivateEngine)
	enginesRoutes.GET("/:id", products.GetEngine)
	enginesRoutes.PATCH("/:id", products.UpdateEngine)
	/*ENGINES ROUTES*/

	/*ACCESSORIES ROUTES*/
	accRoutes := apiv1.Group("/accessories")
	accRoutes.Use(echojwt.WithConfig(configJwt))
	accRoutes.Use(nmiddleware.CheckRoleAndPermissions)

	accRoutes.GET("", products.GetAccessories)
	accRoutes.POST("", products.InsertAccessory)
	accRoutes.DELETE("/:id", products.DeactivateAccessory)
	accRoutes.GET("/:id", products.GetAccessory)
	accRoutes.PATCH("/:id", products.UpdateAccessory)

	//I know, i know.... this should problably a separate group, but you know how things are.....
	accRoutes.GET("/types", products.GetAccessoriesTypes)
	accRoutes.POST("/types", products.InsertAccessoryType)
	accRoutes.DELETE("/types/:id", products.DeactivateAccessoryType)
	accRoutes.PATCH("/types/:id", products.UpdateAccessoryType)
	/*ACCESSORIES ROUTES*/

	/*BOATS ROUTES*/
	boatsRoutes := apiv1.Group("/boats")
	boatsRoutes.Use(echojwt.WithConfig(configJwt))
	boatsRoutes.Use(nmiddleware.CheckRoleAndPermissions)

	boatsRoutes.GET("", products.GetBoats)
	boatsRoutes.POST("", products.InsertBoat)
	boatsRoutes.GET("/:id", products.GetBoat)
	boatsRoutes.PATCH("/:id", products.UpdateBoat)

	boatsRoutes.POST("/:id/ads/:id_mean", products.InsertBoatAd)
	boatsRoutes.POST("/:id/accessories/:id_acc", products.InsertBoatAccessory)
	boatsRoutes.POST("/:id/engines/:id_eng", products.InsertBoatEngine)

	boatsRoutes.DELETE("/:id/accessories/:id_acc", products.RemoveBoatAccessory)
	boatsRoutes.DELETE("/:id/engines/:id_eng", products.RemoveBoatEngine)
	boatsRoutes.DELETE("/:id/ads/:id_mean", products.RemoveBoatAd)

	boatsRoutes.GET("/:id/accessories", products.GetBoatAccessories)
	boatsRoutes.GET("/:id/engines", products.GetBoatEngines)
	boatsRoutes.GET("/:id/ads", products.GetBoatAds)

	/*BOATS ROUTES*/

	/*SALES ROUTES*/
	salRoutes := apiv1.Group("/sales")
	salRoutes.Use(echojwt.WithConfig(configJwt))
	salRoutes.Use(nmiddleware.CheckRoleAndPermissions)

	salRoutes.GET("/communication-means", sales.GetComMeans)
	salRoutes.POST("/communication-means", sales.InsertComMeans)
	salRoutes.DELETE("/communication-means/:id", sales.DeactivateComMeans)
	salRoutes.PATCH("/communication-means/:id", sales.UpdateComMeans)

	salRoutes.GET("/negotiations/:id", sales.GetNegotiation)
	salRoutes.GET("/negotiations", sales.GetNegotiations)
	salRoutes.POST("/negotiations", sales.InsertNegotiation)
	salRoutes.PATCH("/negotiations/:id", sales.UpdateNegotiation)
	salRoutes.PATCH("/negotiations/:id/advance", sales.UpdateNegotiationAdvanceStage)

	salRoutes.POST("/negotiations/:id/history", sales.InsertNegotiationHistory)
	salRoutes.GET("/negotiations/:id/history", sales.GetNegotiationHistory)

	/*SALES ORDERS ROUTES*/
	salOrdRoutes := salRoutes.Group("/orders")
	salOrdRoutes.Use(echojwt.WithConfig(configJwt))
	salOrdRoutes.Use(nmiddleware.CheckRoleAndPermissions)
	salOrdRoutes.POST("/negotiations/history/:id", sales.InsertSalesOrder)
	salOrdRoutes.GET("/:id", sales.GetSalesOrder)
	salOrdRoutes.GET("/:id/itens", sales.GetSalesOrderItens)

	salOrdRoutes.PATCH("/:id/boat/:id_boat", sales.UpdateBoatSalesOrder)
	salOrdRoutes.PATCH("/:id/engine/:id_engine", sales.UpdateEngineSalesOrder)
	salOrdRoutes.PATCH("/:id/accessory/:id_accessory", sales.UpdateAccessorySalesOrder)

	/*SALES ORDERS ROUTES*/

	/*SALES CUSTOMERS ROUTES*/
	salCustRoutes := salRoutes.Group("/customers")
	salCustRoutes.Use(echojwt.WithConfig(configJwt))
	salCustRoutes.Use(nmiddleware.CheckRoleAndPermissions)

	salCustRoutes.GET("", sales.GetCustomers)
	salCustRoutes.GET("/birthdays", sales.GetCustomersBirthday)

	salCustRoutes.GET("/:id", sales.GetCustomer)
	salCustRoutes.PATCH("/:id", sales.UpdateCustomer)
	salCustRoutes.GET("/:id/negotiations/history", sales.GetCustomerNegotiationsHistories)
	/*SALES CUSTOMERS ROUTES*/

	/*SALES ROUTES*/

	/*SALES REPORTS ROUTES*/
	salReps := salRoutes.Group("/reports")
	salReps.Use(echojwt.WithConfig(configJwt))
	salReps.Use(nmiddleware.CheckRoleAndPermissions)
	salReps.GET("/negotiations", reports.GetNegotiationsReport)
	salReps.GET("/sales-orders", reports.GetSalesOrdersReport)

	/*SALES REPORTS ROUTES*/

	e.Logger.Fatal(e.Start(":8080"))
}
