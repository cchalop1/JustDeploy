package application

import (
	"errors"
	"path/filepath"
	"strings"

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

	service := domain.Service{
		Id:          utils.GenerateRandomPassword(5),
		Status:      "ready_to_deploy",
		Name:        Name,
		CurrentPath: repoPath,
		Type:        "github_repo",
		ImageName:   Name,
	}

	deployService.DatabaseAdapter.SaveService(service)

	return service, nil
}
