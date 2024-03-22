package application

import (
	"cchalop1.com/deploy/internal/api/service"
)

func RemoveApplicationById(deployService *service.DeployService, deployId string) error {
	deploy, err := deployService.DatabaseAdapter.GetDeployById(deployId)
	if err != nil {
		return err
	}
	// TODO: connect to docker client
	deployService.DockerAdapter.Delete(deploy.Name, true)
	deployService.DatabaseAdapter.DeleteDeploy(deploy)
	return nil
}
