package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func ReDeployAppHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		containerName := c.Param("name")

		application.RemoveApplication(containerName, deployService)
		application.DeployApplication(deployService)

		return c.JSON(http.StatusOK, dto.ResponseApi{Message: "Application is redeploy"})
	}
}
