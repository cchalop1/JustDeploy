package application

import (
	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api/dto"
)

func GetDeployConfig(databaseAdapter *adapter.DatabaseAdapter, filesystemAdapter *adapter.FilesystemAdapter) dto.DeployConfigDto {
	configDeploy := databaseAdapter.GetState()

	currentPath, err := filesystemAdapter.GetCurrentPath()
	if err == nil && configDeploy.AppConfig.PathToSource == "" {
		configDeploy.AppConfig.PathToSource = currentPath

	}
	return configDeploy
}
