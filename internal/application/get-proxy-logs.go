package application

import (
	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
)

func GetServerProxyLogs(deployService *service.DeployService, serverId string) ([]dto.Logs, error) {
	server, err := deployService.DatabaseAdapter.GetServerById(serverId)

	logs := []dto.Logs{}

	if err != nil {
		return nil, err
	}

	// TODO: move all the connect client to a midleware
	deployService.DockerAdapter.ConnectClient(server)

	dockerLogs, err := deployService.DockerAdapter.GetLogsOfContainer(adapter.ROUTER_NAME)

	if err != nil {
		return nil, err
	}

	for _, log := range dockerLogs {

		if len(log) < 30 {
			continue
		}

		datePart := log[0:30]
		messagePart := log[30:]
		logs = append(logs, dto.Logs{
			Date:    datePart,
			Message: messagePart,
		})
	}

	return logs, nil
}
