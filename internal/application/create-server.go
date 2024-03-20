package application

import (
	"strconv"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/internal/utils"
)

func CreateServer(deployService *service.DeployService, createNewServer dto.ConnectNewServerDto) bool {
	// serverCount := deployService.DatabaseAdapter.
	// TODO: change to get from the database
	serverCount := 1
	Name := "Server " + strconv.Itoa(serverCount)
	server := domain.Server{
		Id:          utils.GenerateRandomPassword(5),
		Name:        Name,
		Domain:      createNewServer.Domain,
		Password:    createNewServer.Password,
		SshKey:      createNewServer.SshKey,
		Ip:          "",
		CreatedDate: "",
		Status:      "Installing",
	}

	deployService.DatabaseAdapter.SaveServer(server)

	// TODO: try to make one ssh connection
	// TODO: make async call to setup server

	return true
}
