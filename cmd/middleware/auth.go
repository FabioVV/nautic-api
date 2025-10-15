package middleware

import (
	"net/http"
	"slices"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"nautic/auth"
)

var RoutesPermissions = map[string]string{
	// "POST:/api/v1/users":       "users:create",
	// "GET:/api/v1/users/:id":    "users:view",
	// "PATCH:/api/v1/users/:id":  "users:update",
	// "DELETE:/api/v1/users/:id": "users:delete",

	//sistema
	"GET:/api/v1/users": "users:view",
	"GET:/api/v1/roles": "roles:view",

	//vendas
	"GET:/api/v1/sales/oportunities":        "sales_oportunities:view",
	"GET:/api/v1/sales/customers":           "sales_customers:view",
	"GET:/api/v1/sales/negotiations":        "negotiation_panel:view",
	"GET:/api/v1/sales/communication-means": "communication_means:view",

	//produtos
	"GET:/api/v1/accessories/types": "accessories_types:view",
	"GET:/api/v1/accessories":       "accessories:view",
	"GET:/api/v1/boats":             "pboats:view",
	"GET:/api/v1/engines":           "engines:view",

	//relatorios
	"GET:/api/v1/reports/negotiations":      "reports_negotiations:view",
	"GET:/api/v1/reports/sales-orders":      "reports_sales_orders:view",
	"GET:/api/v1/reports/lost-negotiations": "reports_lost_negotiations:view",
}

func CheckRoleAndPermissions(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse user credentials"})
		}
		claims, ok := user.Claims.(*auth.JwtCustomClaims)
		if !ok {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse user credentials claims"})
		}

		if len(claims.Roles) > 0 && claims.Roles[0] == "Admin" {
			return next(c)
		}

		routePermissionKey := c.Request().Method + ":" + c.Path()
		routePermission := RoutesPermissions[routePermissionKey]

		if slices.Contains(claims.Permissions, routePermission) {
			return next(c)
		}

		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "User does not have permission for the requested resource"})
	}
}
