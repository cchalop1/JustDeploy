package application

import (
	"cchalop1.com/deploy/internal/api/service"
)

func RemoveApplicationById(deployService *service.DeployService, deployId string) error {
	// deploy, err := deployService.DatabaseAdapter.GetDeployById(deployId)
	// if err != nil {
	// 	return err
	// }
	// server, err := deployService.DatabaseAdapter.GetServerById(deploy.ServerId)
	// if err != nil {
	// 	return err
	// }
	// deployService.DockerAdapter.ConnectClient(server)
	// services := deployService.DatabaseAdapter.GetServicesByDeployId(deployId)

	// if len(services) > 0 {
	// 	for _, service := range services {
	// 		DeleteService(deployService, service.Id)
	// 	}
	// }

	// deployService.DockerAdapter.Delete(deploy.GetDockerName())
	// deployService.DatabaseAdapter.DeleteDeploy(deploy)
	// deployService.DatabaseAdapter.DeleteLogFile(deployId)

	return nil
}
