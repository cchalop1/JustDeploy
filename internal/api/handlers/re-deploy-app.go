package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/usecase"
	"cchalop1.com/deploy/internal/application"
	"cchalop1.com/deploy/internal/domain"
	"github.com/labstack/echo/v4"
)

func ReDeployAppHandler(deployUseCase *usecase.DeployUseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		containerName := c.Param("name")

		application.RemoveApplication(containerName, deployUseCase)
		application.DeployApplication(deployUseCase)

		return c.JSON(http.StatusOK, domain.ResponseApi{Message: "Application is redeploy"})
	}
}
