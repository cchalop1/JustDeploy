package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func CreateDeployementHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		postCreateDeploymentRequest := dto.AppConfigDto{}
		err := c.Bind(&postCreateDeploymentRequest)
		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		// TODO: create a models for appConfig
		deployService.DeployConfig.AppConfig = postCreateDeploymentRequest

		application.DeployApplication(deployService)

		return c.JSON(http.StatusOK, dto.ResponseApi{Message: "Application is deploy"})
	}
}
