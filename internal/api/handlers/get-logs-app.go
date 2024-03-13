package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func GetLogsHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		containerName := c.Param("name")
		logs := application.GetApplicationLogs(deployService.DockerAdapter, containerName)

		return c.JSON(http.StatusOK, logs)
	}
}
