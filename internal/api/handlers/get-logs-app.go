package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/usecase"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func GetLogsHandler(deployUseCase *usecase.DeployUseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		containerName := c.Param("name")
		logs := application.GetApplicationLogs(deployUseCase.DockerAdapter, containerName)

		return c.JSON(http.StatusOK, logs)
	}
}
