package application

import (
	"cchalop1.com/deploy/internal/api/service"
)

func RemoveApplicationById(deployService *service.DeployService, deployId string) error {
	deploy, err := deployService.DatabaseAdapter.GetDeployById(deployId)
	if err != nil {
		return err
	}
	server, err := deployService.DatabaseAdapter.GetServerById(deploy.ServerId)
	if err != nil {
		return err
	}
	// TODO: remove services on deploy if there is
	deployService.DockerAdapter.ConnectClient(server)
	deployService.DockerAdapter.Delete(deploy.GetDockerName(), false)
	deployService.DatabaseAdapter.DeleteDeploy(deploy)
	return nil
}
