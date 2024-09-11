package application

import (
	"cchalop1.com/deploy/internal/api/graph/model"
	"cchalop1.com/deploy/internal/api/service"
)

func GetProjectById(deployService *service.DeployService, id string) (*model.Project, error) {
	product, err := deployService.DatabaseAdapter.GetProjectById(id)

	if err != nil {
		return nil, err
	}

	return &model.Project{
		ID:       product.Id,
		Name:     product.Name,
		Path:     product.Path,
		Services: []*model.Service{},
		Apps:     []*model.App{},
	}, nil
}
