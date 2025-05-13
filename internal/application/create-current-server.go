package application

import (
	"time"

	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/internal/utils"
)

func CreateCurrentServer(deployService *service.DeployService, port string) (domain.Server, error) {
	currentIp, err := deployService.NetworkAdapter.GetCurrentIP()

	if err != nil {
		return domain.Server{}, err
	}

	server := deployService.DatabaseAdapter.GetServer()

	if server.Id != "" {
		return server, nil
	}

	server = domain.Server{
		Id:          utils.GenerateRandomPassword(5),
		Name:        "Local Server",
		Ip:          currentIp,
		Port:        port,
		CreatedDate: time.Now(),
		Status:      "Installing",
	}

	deployService.DatabaseAdapter.SaveServer(server)

	return server, nil
}
