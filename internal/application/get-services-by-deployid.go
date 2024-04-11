package application

import (
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
)

func GetServicesByDeployId(deployService *service.DeployService, deployId string) []domain.Service {
	return deployService.DatabaseAdapter.GetServiceByDeployId(deployId)
}
