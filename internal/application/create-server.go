package application

import (
	"errors"
	"strconv"
	"time"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/internal/utils"
)

func CreateServer(deployService *service.DeployService, createNewServer dto.ConnectNewServerDto) (string, error) {
	serverCount := deployService.DatabaseAdapter.CountServer() + 1
	Name := "Server " + strconv.Itoa(serverCount)

	server := domain.Server{
		Id:          utils.GenerateRandomPassword(5),
		Name:        Name,
		Domain:      createNewServer.Domain,
		Password:    createNewServer.Password,
		SshKey:      createNewServer.SshKey,
		Ip:          createNewServer.Ip,
		CreatedDate: time.Now(),
		Status:      "Installing",
		Email:       createNewServer.Email,
	}

	serverList := deployService.DatabaseAdapter.GetServers()

	for _, s := range serverList {
		if s.Ip == server.Ip {
			return server.Id, errors.New("server already exist")
		}
	}

	projects := deployService.DatabaseAdapter.GetProjects()

	for _, project := range projects {
		if project.ServerId == nil {
			project.ServerId = &server.Id
		}
		deployService.DatabaseAdapter.SaveProject(project)
	}

	deployService.DatabaseAdapter.SaveServer(server)

	ConnectAndSetupServer(deployService, server)

	return server.Id, nil
}
