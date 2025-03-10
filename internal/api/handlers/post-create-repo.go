package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

type CreateRepoRequest struct {
	RepoUrl string `json:"repoUrl"`
}

// ServiceExposeSettings représente les paramètres d'exposition d'un service
type ServiceExposeSettings struct {
	IsExposed  bool   `json:"isExposed"`
	SubDomain  string `json:"subDomain"`
	Tls        bool   `json:"tls"`
	ExposePort string `json:"exposePort"`
}

// ServiceResponse représente la réponse du service pour la documentation Swagger
type ServiceResponse struct {
	Id             string                `json:"id"`
	Type           string                `json:"type"`
	Url            string                `json:"url"`
	Name           string                `json:"name"`
	Envs           []dto.Env             `json:"envs"`
	Status         string                `json:"status"`
	ImageName      string                `json:"imageName"`
	ImageUrl       string                `json:"imageUrl"`
	CurrentPath    string                `json:"currentPath"`
	DockerHubUrl   string                `json:"dockerHubUrl"`
	ExposeSettings ServiceExposeSettings `json:"exposeSettings"`
}

// @Summary Create a new service from a GitHub repository
// @Description Creates a new service by cloning and configuring a GitHub repository
// @Tags repository
// @Accept json
// @Produce json
// @Param request body CreateRepoRequest true "Repository URL information"
// @Success 201 {object} ServiceResponse "Service created successfully"
// @Failure 400 {object} map[string]string "Bad request - Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/repo [post]
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
