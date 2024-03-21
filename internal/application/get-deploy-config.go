package application

import (
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
)

func GetDeployConfig(deployService *service.DeployService) dto.DeployConfigDto {
	deployConfig := dto.DeployConfigDto{}

	currentPath, err := deployService.FilesystemAdapter.GetCurrentPath()

	if err == nil {
		deployConfig.PathToSource = currentPath
	}

	deployConfig.SourceType = "Local Folder"

	return deployConfig
}
