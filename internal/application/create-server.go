package application

import (
	"fmt"
	"strconv"

	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/internal/utils"
)

func CreateServer(deployService *service.DeployService, createNewServer dto.ConnectNewServerDto) bool {
	serverCount := deployService.DatabaseAdapter.CountServer() + 1
	Name := "Server " + strconv.Itoa(serverCount)

	fmt.Println("Creating new server")

	server := domain.Server{
		Id:          utils.GenerateRandomPassword(5),
		Name:        Name,
		Domain:      createNewServer.Domain,
		Password:    createNewServer.Password,
		SshKey:      createNewServer.SshKey,
		Ip:          createNewServer.Ip,
		CreatedDate: "",
		Status:      "Installing",
	}

	deployService.DatabaseAdapter.SaveServer(server)

	sshAdapter := adapter.NewSshAdapter()

	fmt.Println("connect tot server")
	err := sshAdapter.Connect(dto.ConnectNewServerDto{
		Domain:   server.Domain,
		SshKey:   server.SshKey,
		Password: server.Password,
		User:     "root",
	})

	if err != nil {
		fmt.Println("Error connecting to server", err)
		return false
	}

	ConnectAndSetupServer(deployService, server)

	return true
}
