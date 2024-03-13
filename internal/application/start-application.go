package application

import "cchalop1.com/deploy/internal/api/service"

func StartApplication(deployService *service.DeployService, containerName string) {
	deployService.DockerAdapter.Start(containerName)
	deployService.DeployConfig.AppStatus = "Stoped"
}
