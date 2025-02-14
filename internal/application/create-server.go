package application

import (
	"time"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/internal/utils"
)

func CreateServer(deployService *service.DeployService, createNewServer dto.ConnectNewServerDto) (string, error) {
	server := domain.Server{
		Id:          utils.GenerateRandomPassword(5),
		Name:        "Server 1",
		Domain:      createNewServer.Domain,
		Password:    createNewServer.Password,
		SshKey:      createNewServer.SshKey,
		Ip:          createNewServer.Ip,
		CreatedDate: time.Now(),
		Status:      "Installing",
		Email:       createNewServer.Email,
	}

	deployService.DatabaseAdapter.SaveServer(server)

	ConnectAndSetupServer(deployService, server)

	return server.Id, nil
}
