package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

type CreateRepoRequest struct {
	RepoUrl string `json:"repoUrl"`
}

func PostCreateRepoHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req CreateRepoRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid request format",
			})
		}

		// Validate required fields
		if req.RepoUrl == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "RepoUrl is required",
			})
		}

		service, err := application.CreateServiceFromGithubRepo(deployService, req.RepoUrl)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusCreated, service)
	}
}
