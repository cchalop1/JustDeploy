package application

import (
	"errors"

	"cchalop1.com/deploy/internal/api/service"
)

func ReDeployApplication(deployService *service.DeployService, serviceName string) error {
	var err error

	settings := deployService.DatabaseAdapter.GetSettings()
	if settings.GithubToken == "" {
		return errors.New("GitHub token not found. Please configure GitHub integration in settings")
	}

	services := GetServices(deployService)

	for _, s := range services {
		if s.Name == serviceName {
			err = deployService.GitAdapter.CloneRepository(s.RepoUrl, s.CurrentPath, settings.GithubToken)
			if err != nil {
				return err
			}
			break
		}
	}

	err = DeployApplication(deployService)
	return err
}
