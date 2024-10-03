package application

import (
	"fmt"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
)

func GetProjectSettings(deployService *service.DeployService, projectId string) (dto.ProjectSettingsDto, error) {
	project, err := deployService.DatabaseAdapter.GetProjectById(projectId)

	if err != nil {
		return dto.ProjectSettingsDto{}, err
	}

	folders, err := deployService.FilesystemAdapter.GetFolders(project.Path)
	fmt.Println(folders)

	// folders = append(folders, dto.PathDto{
	// 	Name:     deployService.FilesystemAdapter.GetFolderName(project.Path),
	// 	FullPath: project.Path,
	// })

	if err != nil {
		return dto.ProjectSettingsDto{}, err
	}

	foldersList := []dto.PathDto{}

	for _, folder := range folders {
		keepFolder := true
		for _, service := range project.Services {
			if service.CurrentPath == folder.FullPath {
				keepFolder = false
			}
		}
		if keepFolder {
			foldersList = append(foldersList, folder)
		}
	}

	projectSetting := dto.ProjectSettingsDto{
		CurrentPath:       project.Path,
		CurrentFolderName: deployService.FilesystemAdapter.GetFolderName(project.Path),
		Folders:           foldersList,
	}

	return projectSetting, nil
}
