package application

// func GetApplicationLogs(deployService *service.DeployService, deployId string) ([]domain.Logs, error) {
// 	deploy, err := deployService.DatabaseAdapter.GetDeployById(deployId)
// 	logs := []domain.Logs{}

// 	if err != nil {
// 		return logs, err
// 	}

// 	server, err := deployService.DatabaseAdapter.GetServerById(deploy.ServerId)

// 	if err != nil {
// 		return logs, err
// 	}

// 	// TODO: move all the connect client to a midleware
// 	deployService.DockerAdapter.ConnectClient(server)

// 	dockerLogs, err := deployService.DockerAdapter.GetLogsOfContainer(deploy.GetDockerName())
// 	if err != nil {
// 		return deployService.DatabaseAdapter.GetLogs(deployId)
// 	}

// 	for _, log := range dockerLogs {

// 		if len(log) < 30 {
// 			continue
// 		}

// 		datePart := log[0:30]
// 		messagePart := log[30:]
// 		logs = append(logs, domain.Logs{
// 			Date:    datePart,
// 			Message: messagePart,
// 		})
// 	}

// 	err = deployService.DatabaseAdapter.SaveLogs(deployId, logs)
// 	if err != nil {
// 		return logs, err
// 	}

// 	return logs, nil
// }
