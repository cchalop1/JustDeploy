package middleware

import (
	"net/http"
	"strings"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"github.com/labstack/echo/v4"
)

// APIKeyAuth middleware validates the API key in request headers or query parameters
func APIKeyAuth(deployService *service.DeployService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path

			if strings.HasPrefix(path, "/github/auth") ||
				!strings.HasPrefix(path, "/api") ||
				strings.HasPrefix(path, "/api/github/connect") {
				return next(c)
			}

			// Get settings from database
			settings := deployService.DatabaseAdapter.GetSettings()

			// Skip authentication if API key is not set yet
			if settings.ApiKey == "" {
				return next(c)
			}

			// Check for API key in header
			apiKey := c.Request().Header.Get("X-API-Key")

			// If not in header, check query parameter
			if apiKey == "" {
				apiKey = c.QueryParam("api_key")
			}

			// If API key is not provided or doesn't match
			if apiKey == "" || apiKey != settings.ApiKey {
				return c.JSON(http.StatusUnauthorized, dto.ResponseApi{
					Message: "Invalid or missing API key",
				})
			}

			// API key is valid, continue with the request
			return next(c)
		}
	}
}
