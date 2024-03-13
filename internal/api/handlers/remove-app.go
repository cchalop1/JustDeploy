package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/usecase"
	"cchalop1.com/deploy/internal/application"
	"cchalop1.com/deploy/internal/domain"
	"github.com/labstack/echo/v4"
)

func RemoveApplicationHandler(deployUseCase *usecase.DeployUseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		applicationName := c.Param("name")

		application.RemoveApplication(applicationName, deployUseCase)

		// h.databaseAdapter.SaveState(h.deployConfig)
		return c.JSON(http.StatusOK, domain.ResponseApi{Message: "Application is removed"})
	}
}
