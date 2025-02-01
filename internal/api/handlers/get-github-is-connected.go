package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func GetGithubIsConnectedHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		isConnected := application.GithubIsConnected(deployService)
		return c.JSON(http.StatusOK, dto.GithubIsConnected{IsConnected: isConnected})
	}
}
