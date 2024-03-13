package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func ConnectAndSetupServerHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		connectServerDto := dto.ConnectServerDto{}

		err := c.Bind(&connectServerDto)
		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		deployService.DeployConfig.ServerConfig = connectServerDto
		deployService.DockerAdapter = application.ConnectAndSetupServer(deployService)
		deployService.DeployConfig.DeployStatus = "appconfig"

		return c.JSON(http.StatusOK, deployService.DeployConfig)
	}
}
