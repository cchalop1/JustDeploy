package application

import (
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/internal/utils"
)

func CreateApp(deployService *service.DeployService, createApp dto.CreateAppDto) (domain.App, error) {
	project, err := deployService.DatabaseAdapter.GetProjectById(createApp.ProjectId)

	if err != nil {
		return domain.App{}, err
	}

	app := domain.App{
		Id:   utils.GenerateRandomPassword(8),
		Path: createApp.Path,
		Name: deployService.FilesystemAdapter.GetFolderName(createApp.Path),

		// TODO: get there information from the project
		IsDockerFile:    false,
		IsDockerCompose: false,
	}

	project.Apps = append(project.Apps, app)

	err = deployService.DatabaseAdapter.SaveProject(*project)

	if err != nil {
		return domain.App{}, err
	}

	return app, nil
}
