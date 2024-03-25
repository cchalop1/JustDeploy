package application

import "cchalop1.com/deploy/internal/api/service"

func GetApplicationLogs(deployService *service.DeployService, deployId string) ([]string, error) {
	deploy, err := deployService.DatabaseAdapter.GetDeployById(deployId)
	if err != nil {
		return []string{}, err
	}

	server, err := deployService.DatabaseAdapter.GetServerById(deploy.ServerId)
	if err != nil {
		return []string{}, err
	}
	deployService.DockerAdapter.ConnectClient(server)

	return deployService.DockerAdapter.GetLogsOfContainer(deploy.GetDockerName()), nil
}
