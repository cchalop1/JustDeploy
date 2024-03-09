package application

import (
	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/domain"
)

func GetFormDetails(filesystemAdapter *adapter.FilesystemAdapter) domain.DeployConfigDto {
	formDetailsResponse := domain.DeployConfigDto{
		PathToProject:    "",
		DockerFileValid:  false,
		DeployFromStatus: "serverconfig",
		ServerConfig:     domain.ConnectServerDto{},
		AppConfig:        domain.AppConfigDto{},
	}
	// pathToProject := ""

	// flag.StringVar(&pathToProject, "path", "", "Path to the deployment directory")

	// flag.Parse()

	// if pathToProject == "" {
	// 	pathToProject = filesystemAdapter.GetCurrentPath()
	// }

	// formDetailsResponse.AppConfig.Name = filesystemAdapter.GetFolderName(pathToProject)
	// formDetailsResponse.PathToProject = filesystemAdapter.GetFullPathToProject(pathToProject)
	// formDetailsResponse.DockerFileValid = filesystemAdapter.IsWhereIsADockerFileInTheFolder(formDetailsResponse.PathToProject)

	return formDetailsResponse
}
