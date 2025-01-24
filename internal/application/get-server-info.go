package application

import (
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
)

func GetServerInfo(deployService *service.DeployService) (domain.Server, error) {
	return deployService.DatabaseAdapter.GetCurrentServer()
}
