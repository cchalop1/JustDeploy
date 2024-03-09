package application

import (
	"path/filepath"

	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/domain"
)

type DeploymentService struct {
	dockerAdapter *adapter.DockerAdapter
}

func NewDeploymentService(dockerAdapter *adapter.DockerAdapter) *DeploymentService {
	return &DeploymentService{
		dockerAdapter: dockerAdapter,
	}
}

func (s *DeploymentService) DeployApplication(deployConfig domain.DeployConfigDto) error {

	pathToDir, err := filepath.Abs(deployConfig.PathToProject)

	if err != nil {
		return err
	}

	s.dockerAdapter.BuildImage(deployConfig.AppConfig.Name, pathToDir)
	s.dockerAdapter.PullTreafikImage()
	s.dockerAdapter.RunRouter()
	s.dockerAdapter.RunImage(deployConfig)

	return nil
}

func (s *DeploymentService) RemoveApplication(applicationName string) error {
	s.dockerAdapter.Remove(applicationName)
	return nil
}
