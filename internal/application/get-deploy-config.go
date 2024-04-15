package application

import (
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
)

func GetDeployConfig(deployService *service.DeployService, deployId string) dto.DeployConfigDto {
	deployConfig := dto.DeployConfigDto{}

	currentPath, err := deployService.FilesystemAdapter.GetCurrentPath()

	if deployId != "" {
		deploy, err := deployService.DatabaseAdapter.GetDeployById(deployId)

		if err != nil {
			return dto.DeployConfigDto{}
		}
		currentPath = deploy.PathToSource
	}

	dockerfileIsFound := deployService.FilesystemAdapter.FindDockerFile(currentPath)
	dockercomposeIsFound := deployService.FilesystemAdapter.FindDockerComposeFile(currentPath)
	envs := deployService.FilesystemAdapter.LoadEnvsFromFileSystem(currentPath)

	deployConfig.DockerFileFound = dockerfileIsFound
	deployConfig.ComposeFileFound = dockercomposeIsFound

	if len(envs) > 0 {
		deployConfig.EnvFileFound = true
	} else {
		deployConfig.EnvFileFound = false
	}

	deployConfig.Envs = envs

	if err == nil {
		deployConfig.PathToSource = currentPath
	}

	deployConfig.SourceType = "Local Folder"

	return deployConfig
}
