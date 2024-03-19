package application

import (
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/internal/utils"
)

func CreateServer(deployService *service.DeployService, createNewServer dto.ConnectNewServerDto) bool {
	// serverCount := deployService.DatabaseAdapter.
	// TODO: change to get from the database
	serverCount := 1
	Name := "Server " + string(serverCount)
	server := domain.Server{
		// TODO: change to generate a uuid
		Id:          utils.GenerateRandomPassword(5),
		Name:        Name,
		Domain:      createNewServer.Domain,
		Password:    createNewServer.Password,
		SshKey:      createNewServer.SshKey,
		Ip:          "",
		CreatedDate: "",
	}

	deployService.DatabaseAdapter.SaveServer(server)
	// TODO: try to connect to your server with

	return true
}
