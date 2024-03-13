package application

import "cchalop1.com/deploy/internal/api/usecase"

func StopApplication(deployUseCase *usecase.DeployUseCase, containerName string) {
	deployUseCase.DockerAdapter.Stop(containerName)
	deployUseCase.DeployConfig.AppStatus = "Stopped"
}
