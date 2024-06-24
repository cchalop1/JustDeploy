package application

import (
	"fmt"
	"strconv"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/internal/utils"
)

func CreateServer(deployService *service.DeployService, createNewServer dto.ConnectNewServerDto) domain.Server {
	serverCount := deployService.DatabaseAdapter.CountServer() + 1
	Name := "Server " + strconv.Itoa(serverCount)

	fmt.Println("Creating new server")

	server := domain.Server{
		Id:          utils.GenerateRandomPassword(5),
		Name:        Name,
		Domain:      "",
		Password:    createNewServer.Password,
		SshKey:      createNewServer.SshKey,
		Ip:          createNewServer.Ip,
		CreatedDate: "",
		Status:      "Installing",
	}

	deployService.DatabaseAdapter.SaveServer(server)

	// go ConnectAndSetupServer(deployService, server)

	return server
}
