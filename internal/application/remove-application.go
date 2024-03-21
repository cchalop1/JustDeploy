package application

import (
	"cchalop1.com/deploy/internal/api/service"
)

func RemoveApplication(applicationName string, deployService *service.DeployService) error {
	deployService.DockerAdapter.Delete(applicationName, true)
	// deployService.DeployConfig.DeployStatus = "appconfig"
	return nil
}
