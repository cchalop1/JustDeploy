package handlers

import (
	"cchalop1.com/deploy/internal/api/service"
	"github.com/labstack/echo/v4"
)

func PostCreateDatabaseHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}
