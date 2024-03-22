package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func GetDeployByServerIdHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		serverId := c.Param("id")
		deployList := application.GetDeployByServerId(deployService, serverId)
		return c.JSON(http.StatusOK, deployList)
	}
}
