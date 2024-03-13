package application

import (
	"path/filepath"

	"cchalop1.com/deploy/internal/api/usecase"
)

func DeployApplication(deployUseCase *usecase.DeployUseCase) error {
	pathToDir, err := filepath.Abs(deployUseCase.DeployConfig.PathToProject)

	if err != nil {
		return err
	}

	deployUseCase.DockerAdapter.BuildImage(deployUseCase.DeployConfig.AppConfig.Name, pathToDir)
	deployUseCase.DockerAdapter.PullTreafikImage()
	deployUseCase.DockerAdapter.RunRouter()
	deployUseCase.DockerAdapter.RunImage(*deployUseCase.DeployConfig)

	deployUseCase.DeployConfig.DeployStatus = "deployapp"

	if deployUseCase.DeployConfig.AppConfig.EnableTls {
		deployUseCase.DeployConfig.Url = "https://" + deployUseCase.DeployConfig.ServerConfig.Domain
	} else {
		deployUseCase.DeployConfig.Url = "http://" + deployUseCase.DeployConfig.ServerConfig.Domain
	}

	deployUseCase.DeployConfig.AppStatus = "Runing"
	return nil
}
