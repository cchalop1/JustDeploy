package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func RemoveServerHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		serverId := c.Param("id")

		application.RemoveServerById(deployService, serverId)

		return c.JSON(http.StatusOK, dto.ResponseApi{Message: "Server is removed"})
	}
}
