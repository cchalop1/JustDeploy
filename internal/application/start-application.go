package application

import "cchalop1.com/deploy/internal/api/service"

func StartApplication(deployService *service.DeployService, deployId string) error {
	deploy, err := deployService.DatabaseAdapter.GetDeployById(deployId)
	if err != nil {
		return err
	}
	server, err := deployService.DatabaseAdapter.GetServerById(deploy.ServerId)
	if err != nil {
		return err
	}
	deployService.DockerAdapter.ConnectClient(server)
	deployService.DockerAdapter.Start(deploy.GetDockerName())
	deploy.Status = "Runing"
	deployService.DatabaseAdapter.UpdateDeploy(deploy)
	return nil
}
