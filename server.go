package main

import (
	"log"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"nautic/auth"
	"nautic/cmd/handlers/audit"
	auth_h "nautic/cmd/handlers/auth"
	"nautic/cmd/handlers/products"
	"nautic/cmd/handlers/reports"
	"nautic/cmd/handlers/roles"
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
	userRoutes.GET("/:id/permissions", users.GetUserPermissions)

	userRoutes.PATCH("/:id/permissions/:id_perm", users.InsertUserPermission)
	userRoutes.DELETE("/:id/permissions/:id_perm", users.RemoveUserPermission)

	userRoutes.PATCH("/:id", users.UpdateUser)
	userRoutes.DELETE("/:id", users.DeactivateUser)
	/*USER ROUTES*/

	/*AUDIT ROUTES*/
	auditRoutes := apiv1.Group("/audit")
	auditRoutes.Use(echojwt.WithConfig(configJwt))
	auditRoutes.Use(nmiddleware.CheckRoleAndPermissions)

	auditRoutes.GET("/logs", audit.GetLogs)

	/*AUDIT ROUTES*/

	/*ROLES ROUTES*/
	rolesRoutes := apiv1.Group("/roles")
	rolesRoutes.Use(echojwt.WithConfig(configJwt))
	rolesRoutes.Use(nmiddleware.CheckRoleAndPermissions)

	rolesRoutes.POST("", roles.InsertRole)
	rolesRoutes.GET("", roles.GetRoles)
	rolesRoutes.GET("/:id", roles.GetRole)
	rolesRoutes.GET("/:id/permissions", roles.GetRolePermissions)

	rolesRoutes.PATCH("/:id/permissions/:id_perm", roles.InsertRolePermission)
	rolesRoutes.DELETE("/:id/permissions/:id_perm", roles.RemoveRolePermission)

	rolesRoutes.DELETE("/:id", roles.DeleteRole)
	rolesRoutes.PATCH("/:id", roles.UpdateRole)

	/*ROLES ROUTES*/

	/*PERMISSION ROUTES*/
	permsRoutes := apiv1.Group("/permissions")
	permsRoutes.Use(echojwt.WithConfig(configJwt))
	permsRoutes.Use(nmiddleware.CheckRoleAndPermissions)

	permsRoutes.GET("", auth_h.GetPermissions)
	/*PERMISSION ROUTES*/

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
	boatsRoutes.GET("/:id/files", products.GetBoatFiles)
	boatsRoutes.DELETE("/:id/files/:id_file", products.RemoveBoatFile)

	boatsRoutes.PATCH("/:id", products.UpdateBoat)
	boatsRoutes.POST("/:id/boat-file", products.UploadBoatFile)

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
	salRoutes.GET("/negotiations/alerts", sales.GetNegotiationsAlerts)
	salRoutes.POST("/negotiations", sales.InsertNegotiation)
	salRoutes.PATCH("/negotiations/:id", sales.UpdateNegotiation)
	salRoutes.PATCH("/negotiations/:id/deactivate", sales.LostNegotiation)
	salRoutes.PATCH("/negotiations/:id/reactivate", sales.ReactivateNegotiation)
	salRoutes.PATCH("/negotiations/:id/advance", sales.UpdateNegotiationAdvanceStage)

	salRoutes.POST("/negotiations/:id/history", sales.InsertNegotiationHistory)
	salRoutes.GET("/negotiations/:id/history", sales.GetNegotiationHistory)

	/*SALES ORDERS ROUTES*/
	salOrdRoutes := salRoutes.Group("/orders")
	salOrdRoutes.Use(echojwt.WithConfig(configJwt))
	salOrdRoutes.Use(nmiddleware.CheckRoleAndPermissions)
	salOrdRoutes.POST("/negotiations/history/:id", sales.InsertSalesOrder)
	salOrdRoutes.GET("/:id", sales.GetSalesOrder)
	salOrdRoutes.POST("/:id/send-quote", sales.SendQuoteViaEmail)

	/*SALES ORDERS ROUTES QUOTES*/
	quoteRoutes := apiv1.Group("/sales/orders/:id/quote") // NO AUTH
	// permsRoutes.Use(echojwt.WithConfig(configJwt))
	// permsRoutes.Use(nmiddleware.CheckRoleAndPermissions)
	quoteRoutes.GET("", sales.GetSalesOrderQuote)
	quoteRoutes.GET("/itens", sales.GetSalesOrderItens)
	quoteRoutes.GET("/files", sales.GetSalesOrderQuoteBoatFiles)

	/*SALES ORDERS ROUTES QUOTES*/

	salOrdRoutes.DELETE("/:id", sales.CancelSalesOrder)
	salOrdRoutes.DELETE("/:id/accessory/:id_accessory", sales.RemoveAccessorySalesOrder)

	salOrdRoutes.GET("/:id/itens", sales.GetSalesOrderItens)

	salOrdRoutes.POST("/:id/so-files", sales.UploadSalesOrderFile)
	salOrdRoutes.GET("/:id/files", sales.GetSalesOrderFiles)
	salOrdRoutes.DELETE("/:id/files/:id_file", sales.RemoveSalesOrderFile)
	salOrdRoutes.PATCH("/:id/files/:id_file/change-type", sales.ChangeSalesOrderFileType)

	salOrdRoutes.PATCH("/:id/boat/:id_boat", sales.UpdateBoatSalesOrder)
	salOrdRoutes.PATCH("/:id/engine/:id_engine", sales.UpdateEngineSalesOrder)
	salOrdRoutes.PATCH("/:id/accessory/:id_accessory", sales.UpdateAccessorySalesOrder)
	salOrdRoutes.PATCH("/:id/upgrade-quote", sales.UpgradeQuoteToSalesOrder)
	salOrdRoutes.PATCH("/:id/accessory/:id_accessory/change-qty", sales.UpdateAccessoryQtySalesOrder)

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
	salReps.GET("/lost-negotiations", reports.GetLostNegotiationsReport)

	/*SALES REPORTS ROUTES*/
	e.Static("/uploads", "uploads")

	e.Logger.Fatal(e.Start(":8080"))
}
