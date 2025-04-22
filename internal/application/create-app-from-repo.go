package application

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/internal/utils"
)

func CreateServiceFromGithubRepo(deployService *service.DeployService, repoUrl string) (domain.Service, error) {
	// First check if we have a GitHub token
	settings := deployService.DatabaseAdapter.GetSettings()
	if settings.GithubToken == "" {
		return domain.Service{}, errors.New("GitHub token not found. Please configure GitHub integration in settings")
	}

	Name := strings.Split(repoUrl, "/")[1]

	// Clone the repository to a temporary directory
	tempDir := deployService.FilesystemAdapter.GetTempDir()
	repoPath := filepath.Join(tempDir, Name)

	// Use the GitHub token when cloning
	err := deployService.GitAdapter.CloneRepository(repoUrl, repoPath, settings.GithubToken)
	if err != nil {
		return domain.Service{}, err
	}

	// Create the main application service
	service := domain.Service{
		Id:          utils.GenerateRandomPassword(5),
		Status:      "ready_to_deploy",
		Name:        Name,
		CurrentPath: repoPath,
		RepoUrl:     repoUrl,
		Type:        "github_repo",
		ImageName:   Name,
		Envs: []dto.Env{
			dto.Env{
				Name:  "",
				Value: "",
			},
		},
		ExposeSettings: BuildExposeSettings(deployService, Name, true, "80"),
	}

	// Check for docker-compose file
	if deployService.FilesystemAdapter.FindDockerComposeFile(repoPath) {
		// Get services from docker-compose
		services, err := deployService.FilesystemAdapter.GetComposeConfigOfDeploy(repoPath)
		fmt.Println(services)
		if err != nil {
			return domain.Service{}, fmt.Errorf("error reading docker-compose: %w", err)
		}

		// Create services for each non-app service in docker-compose
		for key, value := range services {
			if value.HasBuild() {
				continue
			}
			// Create a new service for each additional service (e.g., database)
			additionalService := domain.Service{
				Id:           utils.GenerateRandomPassword(5),
				Type:         "database",
				Name:         key,
				Status:       "ready_to_deploy",
				ImageName:    value.Image,
				DockerHubUrl: utils.GetDockerHubUrl(value.Image),
				ExposeSettings: domain.ServiceExposeSettings{
					IsExposed:  false,
					ExposePort: value.Ports[0],
				},
			}

			// Convert environment variables
			envVars := value.GetEnvironmentVariables()
			additionalService.Envs = append(additionalService.Envs, envVars...)

			// Save the additional service
			deployService.DatabaseAdapter.SaveService(additionalService)

			// Add connection environment variables to the main app
			service.Envs = append(service.Envs, dto.Env{
				Name:  fmt.Sprintf("%s_HOST", strings.ToUpper(key)),
				Value: fmt.Sprintf("%s-db", strings.ToLower(key)),
			})
		}
	}

	// Save the main application service
	deployService.DatabaseAdapter.SaveService(service)

	return service, nil
}
