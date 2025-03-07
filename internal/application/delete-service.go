package application

import (
	"cchalop1.com/deploy/internal/api/service"
)

func DeleteService(deployService *service.DeployService, serviceId string) error {
	s, err := deployService.DatabaseAdapter.GetServiceById(serviceId)

	if err != nil {
		return err
	}

	server := deployService.DatabaseAdapter.GetServer()

	deployService.DockerAdapter.ConnectClient(server)

	deployService.DockerAdapter.Stop(s.GetDockerName())
	deployService.DockerAdapter.Remove(s.GetDockerName())

	return deployService.DatabaseAdapter.DeleteServiceById(s.Id)
}
