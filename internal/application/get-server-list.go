package application

import (
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
)

func GetServerList(deployService *service.DeployService) []domain.Server {
	// TODO: create a list of server response DTO
	return deployService.DatabaseAdapter.GetServers()
}
