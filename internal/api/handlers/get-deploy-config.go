package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/service"
	"github.com/labstack/echo/v4"
)

func GetDeployConfigHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, "ok")
	}
}
