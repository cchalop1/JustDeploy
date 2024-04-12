package application

import (
	"cchalop1.com/deploy/internal/api/service"
)

func DeleteService(deployService *service.DeployService, serviceId string) error {
	s, err := deployService.DatabaseAdapter.GetServiceById(serviceId)
	if err != nil {
		return err
	}

	deploy, err := deployService.DatabaseAdapter.GetDeployById(s.DeployId)
	if err != nil {
		return err
	}

	server, err := deployService.DatabaseAdapter.GetServerById(deploy.ServerId)
	if err != nil {
		return err
	}
	deployService.DockerAdapter.ConnectClient(server)
	deployService.DockerAdapter.Delete(s.Name, false)
	// TODO: remove envs of the service on the deploymeent

	err = deployService.DatabaseAdapter.DeleteServiceById(serviceId)
	return err
}
