package application

import (
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
)

func GetServices(deployService *service.DeployService) []domain.Service {
	return deployService.DatabaseAdapter.GetServices()
}
