package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func GetServerListHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		serverList := application.GetServerList(deployService)
		return c.JSON(http.StatusOK, serverList)
	}
}
