package application

import (
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
)

func GetDeployConfig(deployService *service.DeployService, paramsDeployConfig dto.ParamsDeployConfigDto) dto.DeployConfigDto {
	deployConfig := dto.DeployConfigDto{}

	currentPath := ""

	if paramsDeployConfig.Path == "" {
		currentPath = deployService.FilesystemAdapter.GetCurrentPath()
	} else {
		currentPath = paramsDeployConfig.Path
	}

	if paramsDeployConfig.DeployId != "" {
		deploy, err := deployService.DatabaseAdapter.GetDeployById(paramsDeployConfig.DeployId)

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

	deployConfig.PathToSource = currentPath

	deployConfig.SourceType = "Local Folder"

	return deployConfig
}
