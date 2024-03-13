package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/usecase"
	"cchalop1.com/deploy/internal/application"
	"cchalop1.com/deploy/models"
	"github.com/labstack/echo/v4"
)

func ConnectAndSetupServerHandler(deployUseCase *usecase.DeployUseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		connectServerDto := models.ConnectServerDto{}

		err := c.Bind(&connectServerDto)
		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		// TODO: move to application
		deployUseCase.DockerAdapter = application.ConnectAndSetupServer(connectServerDto)
		deployUseCase.DeployConfig.ServerConfig = connectServerDto
		deployUseCase.DeployConfig.DeployStatus = "appconfig"

		return c.JSON(http.StatusOK, deployUseCase.DeployConfig)
	}
}
