package application

import (
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
)

func UpdateService(deployService *service.DeployService, serviceToUpdate domain.Service) (domain.Service, error) {
	err := deployService.DatabaseAdapter.SaveService(serviceToUpdate)

	return serviceToUpdate, err
}
