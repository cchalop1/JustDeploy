package application

import (
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
)

func GetDeployConfig(deployService *service.DeployService) dto.DeployConfigDto {
	deployConfig := dto.DeployConfigDto{}

	currentPath, err := deployService.FilesystemAdapter.GetCurrentPath()
	dockerfileIsFound := deployService.FilesystemAdapter.FindDockerFile(currentPath)
	dockercomposeIsFound := deployService.FilesystemAdapter.FindDockerComposeFile(currentPath)
	envs := deployService.FilesystemAdapter.LoadEnvsFromFileSystem(currentPath)

	deployConfig.DockerFileFound = dockerfileIsFound
	deployConfig.ComposeFileFound = dockercomposeIsFound
	deployConfig.Envs = envs

	if err == nil {
		deployConfig.PathToSource = currentPath
	}

	deployConfig.SourceType = "Local Folder"

	return deployConfig
}
