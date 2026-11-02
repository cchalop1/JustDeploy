package application

import (
	"time"

	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/internal/utils"
)

func CreateCurrentServer(deployService *service.DeployService, port string) (domain.Server, error) {
	server := deployService.DatabaseAdapter.GetServer()
	if server.Id != "" {
		return server, nil
	}
	// Public IP discovery is advisory; an unavailable external service must not prevent local setup.
	currentIp, _ := deployService.NetworkAdapter.GetCurrentIP()
	server = domain.Server{
		Id:          utils.GenerateRandomPassword(16),
		Name:        "Local Server",
		Ip:          currentIp,
		Port:        port,
		CreatedDate: time.Now(),
		Status:      "Installing",
	}

	return server, deployService.DatabaseAdapter.SaveServer(server)
}
