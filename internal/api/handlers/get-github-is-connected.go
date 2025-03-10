package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

// @Summary Check GitHub connection status
// @Description Checks if the application is connected to GitHub
// @Tags github
// @Accept json
// @Produce json
// @Success 200 {object} dto.GithubIsConnected "Connection status"
// @Router /api/v1/github/is-connected [get]
func GetGithubIsConnectedHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		isConnected := application.GithubIsConnected(deployService)
		return c.JSON(http.StatusOK, dto.GithubIsConnected{IsConnected: isConnected})
	}
}
