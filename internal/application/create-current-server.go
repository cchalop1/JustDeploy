package application

import (
	"errors"
	"strconv"
	"time"

	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/internal/utils"
)

func CreateCurrentServer(deployService *service.DeployService, port string) (domain.Server, error) {
	serverCount := deployService.DatabaseAdapter.CountServer() + 1
	Name := "Server " + strconv.Itoa(serverCount)

	currentIp, err := deployService.NetworkAdapter.GetCurrentIP()

	if err != nil {
		return domain.Server{}, err
	}

	server := domain.Server{
		Id:          utils.GenerateRandomPassword(5),
		Name:        Name,
		Ip:          currentIp,
		Port:        port,
		CreatedDate: time.Now(),
		Status:      "Installing",
	}

	serverList := deployService.DatabaseAdapter.GetServers()

	for _, s := range serverList {
		if s.Ip == server.Ip {
			return server, errors.New("server already exist")
		}
	}

	deployService.DatabaseAdapter.SaveServer(server)

	return server, nil
}
