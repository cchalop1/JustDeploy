package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func StopAppHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		if err := application.StopApplication(deployService, id); err != nil {
			return operationError(err)
		}
		return c.JSON(http.StatusOK, dto.ResponseApi{Message: "Application stopped"})
	}
}
