package application

import "cchalop1.com/deploy/internal/api/usecase"

func RemoveApplication(applicationName string, deployUseCase *usecase.DeployUseCase) error {
	deployUseCase.DockerAdapter.Delete(applicationName, true)
	deployUseCase.DeployConfig.DeployStatus = "appconfig"
	return nil
}
