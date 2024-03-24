package application

import "cchalop1.com/deploy/internal/api/service"

func CreateDatabaseService(deployService *service.DeployService, deployId string) error {
	deploy, err := deployService.DatabaseAdapter.GetDeployById(deployId)

	if err != nil {
		return err
	}

	server, err := deployService.DatabaseAdapter.GetServerById(deploy.ServerId)
	if err != nil {
		return err
	}

	serviceImage := "redis/redis-stack:latest"

	deployService.DockerAdapter.ConnectClient(server)
	deployService.DockerAdapter.PullService(serviceImage)
	deployService.DockerAdapter.RunRedis()

}
