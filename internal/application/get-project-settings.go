package application

import (
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
)

func GetProjectSettings(deployService *service.DeployService, projectId string) (dto.ProjectSettingsDto, error) {
	project, err := deployService.DatabaseAdapter.GetProjectById(projectId)

	if err != nil {
		return dto.ProjectSettingsDto{}, err
	}

	folders, err := deployService.FilesystemAdapter.GetFolders(project.Path)

	if err != nil {
		return dto.ProjectSettingsDto{}, err
	}

	projectSetting := dto.ProjectSettingsDto{
		CurrentPath:       project.Path,
		CurrentFolderName: deployService.FilesystemAdapter.GetFolderName(project.Path),
		Folders:           folders,
	}

	return projectSetting, nil
}
