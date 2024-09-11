package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/http/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func GetLogsHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		deployId := c.Param("id")
		logs, err := application.GetApplicationLogs(deployService, deployId)

		if err != nil {
			return c.JSON(http.StatusNotFound, dto.ResponseApi{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, logs)
	}
}
