package application

import (
	"path/filepath"

	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api/service"
)

func DeployApplication(deployService *service.DeployService) error {
	pathToDir, err := filepath.Abs(deployService.DeployConfig.AppConfig.PathToSource)

	if err != nil {
		return err
	}

	deployService.DeployConfig.AppConfig.PathToSource = adapter.NewFilesystemAdapter().CleanPath(pathToDir)

	deployService.DockerAdapter.BuildImage(deployService.DeployConfig.AppConfig.Name, deployService.DeployConfig.AppConfig.PathToSource)
	deployService.DockerAdapter.PullTreafikImage()
	deployService.DockerAdapter.RunRouter()
	deployService.DockerAdapter.RunImage(*deployService.DeployConfig)

	deployService.DeployConfig.DeployStatus = "deployapp"

	if deployService.DeployConfig.AppConfig.EnableTls {
		deployService.DeployConfig.Url = "https://" + deployService.DeployConfig.ServerConfig.Domain
	} else {
		deployService.DeployConfig.Url = "http://" + deployService.DeployConfig.ServerConfig.Domain
	}

	deployService.DeployConfig.AppStatus = "Runing"

	deployService.DatabaseAdapter.SaveState(*deployService.DeployConfig)
	return nil
}
