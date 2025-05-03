package middleware

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)

// RequestLogger middleware logs information about each API request
func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Start timer
			start := time.Now()

			// Process request
			err := next(c)

			// Log after request is processed
			method := c.Request().Method
			path := c.Request().URL.Path
			status := c.Response().Status
			latency := time.Since(start)

			// Format and print the log
			fmt.Printf("[%s] %s %s - Status: %d - Time: %v\n",
				time.Now().Format(time.RFC3339),
				method,
				path,
				status,
				latency,
			)

			return err
		}
	}
}
