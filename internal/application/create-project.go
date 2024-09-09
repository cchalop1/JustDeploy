package application

import (
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/internal/utils"
)

func CreateProject(deployService *service.DeployService, createProjectDto dto.CreateProjectDto) (string, error) {
	project := domain.Project{
		Id:       utils.GenerateRandomPassword(8),
		Name:     createProjectDto.Name,
		Path:     createProjectDto.Path,
		Services: []domain.Service{},
		Deploy:   []domain.Deploy{},
	}

	err := deployService.DatabaseAdapter.SaveProject(project)
	if err != nil {
		return "", err
	}

	return project.Id, nil
}
