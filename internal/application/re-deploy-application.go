package application

import (
	"cchalop1.com/deploy/internal/api/service"
)

func ReDeployApplication(deployService *service.DeployService, deployId string) error {
	// deploy, err := deployService.DatabaseAdapter.GetDeployById(deployId)
	// if err != nil {
	// 	return err
	// }

	// server, err := deployService.DatabaseAdapter.GetServerById(deploy.ServerId)
	// if err != nil {
	// 	return err
	// }
	// deployService.DockerAdapter.ConnectClient(server)
	// deployService.DockerAdapter.Delete(deploy.GetDockerName())

	// deploy.Status = "Installing"
	// deployService.DatabaseAdapter.UpdateDeploy(deploy)

	// runApplication(deployService, &deploy, server.Domain)

	// if err != nil {
	// 	return err
	// }

	// deploy.Status = "Runing"
	// deployService.DatabaseAdapter.UpdateDeploy(deploy)

	// fmt.Println("Success to redeploy ", deploy.Name)
	return nil
}
