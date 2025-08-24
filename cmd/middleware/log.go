package middleware

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)


func LogReq(next echo.HandlerFunc) echo.HandlerFunc{
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		end := time.Now()

		fmt.Printf("%s %s %s %d\n", c.Request().Method, c.Request().URL, end.Sub(start), c.Response().Status)
		return err
	}
}
