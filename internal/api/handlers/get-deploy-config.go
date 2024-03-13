package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/usecase"
	"github.com/labstack/echo/v4"
)

func GetDeployConfigHandler(deployUseCase *usecase.DeployUseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, deployUseCase.DeployConfig)
	}
}
