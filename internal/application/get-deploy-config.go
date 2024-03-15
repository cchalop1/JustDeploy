package application

import (
	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api/dto"
)

func GetDeployConfig(databaseAdapter *adapter.DatabaseAdapter, filesystemAdapter *adapter.FilesystemAdapter) dto.DeployConfigDto {
	configDeploy := databaseAdapter.GetState()

	// pathToProject := ""

	// flag.StringVar(&pathToProject, "path", "", "Path to the deployment directory")

	// flag.Parse()

	// if pathToProject == "" {
	// 	pathToProject = filesystemAdapter.GetCurrentPath()
	// }

	currentPath, err := filesystemAdapter.GetCurrentPath()
	if err == nil {
		configDeploy.AppConfig.PathToSource = currentPath

	}
	// formDetailsResponse.AppConfig.Name = filesystemAdapter.GetFolderName(pathToProject)
	// formDetailsResponse.PathToProject = filesystemAdapter.GetFullPathToProject(pathToProject)
	// formDetailsResponse.DockerFileValid = filesystemAdapter.IsWhereIsADockerFileInTheFolder(formDetailsResponse.PathToProject)

	return configDeploy
}
