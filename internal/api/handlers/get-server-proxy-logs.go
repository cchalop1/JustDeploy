package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func GetServerProxyLogs(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		serverId := c.Param("id")
		logs, err := application.GetServerProxyLogs(deployService, serverId)

		if err != nil {
			return c.JSON(http.StatusNotFound, dto.ResponseApi{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, logs)

	}
}
