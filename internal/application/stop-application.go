package application

import "cchalop1.com/deploy/internal/api/service"

func StopApplication(deployService *service.DeployService, containerName string) {
	deployService.DockerAdapter.Stop(containerName)
	deployService.DeployConfig.AppStatus = "Stopped"
}
