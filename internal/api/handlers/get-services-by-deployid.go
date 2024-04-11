package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func GetServicesByDeployIdHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		deployId := c.Param("deployId")

		services := application.GetServicesByDeployId(deployService, deployId)

		return c.JSON(http.StatusOK, services)
	}
}
