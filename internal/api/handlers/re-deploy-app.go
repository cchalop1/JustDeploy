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
		serviceName := c.Param("id")

		err := application.ReDeployApplication(deployService, serviceName)

		if err != nil {
			return operationError(err)
		}

		return c.JSON(http.StatusOK, dto.ResponseApi{Message: "Application is redeploy"})
	}
}
