package application

import "cchalop1.com/deploy/internal/api/service"

func StopApplication(deployService *service.DeployService, deployId string) error {
	deploy, err := deployService.DatabaseAdapter.GetDeployById(deployId)
	if err != nil {
		return err
	}

	server, err := deployService.DatabaseAdapter.GetServerById(deploy.ServerId)
	if err != nil {
		return err
	}
	deployService.DockerAdapter.ConnectClient(server)
	deployService.DockerAdapter.Stop(deploy.Name)
	deploy.Status = "Stopped"
	deployService.DatabaseAdapter.UpdateDeploy(deploy)
	return nil
}
