package handlers

import (
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func PostDeployHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		return application.DeployApplication(deployService)
	}
}
