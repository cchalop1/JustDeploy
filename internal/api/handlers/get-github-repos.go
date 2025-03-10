package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

// @Summary Get GitHub repositories
// @Description Retrieves all GitHub repositories accessible to the application
// @Tags github
// @Accept json
// @Produce json
// @Success 200 {array} interface{} "List of GitHub repositories"
// @Router /api/v1/github/repos [get]
func GetGithubRepos(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		repos := application.GetGithubRepos(deployService)
		return c.JSON(http.StatusOK, repos)
	}
}
