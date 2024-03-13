package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func StartAppHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		containerName := c.Param("name")

		application.StartApplication(deployService, containerName)
		return c.JSON(http.StatusOK, dto.ResponseApi{Message: "Application is stoped"})
	}
}
