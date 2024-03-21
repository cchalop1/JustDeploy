package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func RemoveApplicationHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		deployId := c.Param("id")

		application.RemoveApplicationById(deployService, deployId)

		// h.databaseAdapter.SaveState(h.deployConfig)
		return c.JSON(http.StatusOK, dto.ResponseApi{Message: "Application is removed"})
	}
}
