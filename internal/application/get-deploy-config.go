package application

import (
	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api/dto"
)

func GetDeployConfig(databaseAdapter *adapter.DatabaseAdapter) dto.DeployConfigDto {
	configDeploy := databaseAdapter.GetState()

	// pathToProject := ""

	// flag.StringVar(&pathToProject, "path", "", "Path to the deployment directory")

	// flag.Parse()

	// if pathToProject == "" {
	// 	pathToProject = filesystemAdapter.GetCurrentPath()
	// }

	// formDetailsResponse.AppConfig.Name = filesystemAdapter.GetFolderName(pathToProject)
	// formDetailsResponse.PathToProject = filesystemAdapter.GetFullPathToProject(pathToProject)
	// formDetailsResponse.DockerFileValid = filesystemAdapter.IsWhereIsADockerFileInTheFolder(formDetailsResponse.PathToProject)

	return configDeploy
}
