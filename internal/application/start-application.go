package application

import "cchalop1.com/deploy/internal/api/usecase"

func StartApplication(deployUseCase *usecase.DeployUseCase, containerName string) {
	deployUseCase.DockerAdapter.Start(containerName)
	deployUseCase.DeployConfig.AppStatus = "Stoped"
}
