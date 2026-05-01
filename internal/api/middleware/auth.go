package middleware

import (
	"net/http"
	"strings"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func JWTAuth(deployService *service.DeployService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, dto.ResponseApi{
					Message: "Missing or invalid authorization header",
				})
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			settings := deployService.DatabaseAdapter.GetSettings()

			_, err := application.ValidateJWT(tokenStr, settings.JwtSecret)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, dto.ResponseApi{
					Message: "Invalid or expired token",
				})
			}

			return next(c)
		}
	}
}
