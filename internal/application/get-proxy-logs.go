package application

import (
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
)

func GetServerProxyLogs(deployService *service.DeployService, serverId string) ([]domain.Logs, error) {
	logs := []domain.Logs{}

	// dockerLogs, err := deployService.DockerAdapter.GetLogsOfContainer(adapter.ROUTER_NAME)

	// if err != nil {
	// 	return nil, err
	// }

	// for _, log := range dockerLogs {

	// 	if len(log) < 30 {
	// 		continue
	// 	}

	// 	datePart := log[0:30]
	// 	messagePart := log[30:]
	// 	logs = append(logs, domain.Logs{
	// 		Date:    datePart,
	// 		Message: messagePart,
	// 	})
	// }

	return logs, nil
}
