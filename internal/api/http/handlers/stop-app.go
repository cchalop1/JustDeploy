package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/http/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func StopAppHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		application.StopApplication(deployService, id)
		return c.JSON(http.StatusOK, dto.ResponseApi{Message: "Application is stoped"})
	}
}
