package application

import (
	"errors"

	"cchalop1.com/deploy/internal/api/service"
)

func RemoveServerById(deployService *service.DeployService, serverId string) error {
	server, err := deployService.DatabaseAdapter.GetServerById(serverId)
	if err != nil {
		return err
	}

	deployList := deployService.DatabaseAdapter.GetDeployByServerId(serverId)

	if len(deployList) > 0 {
		return errors.New("you can't remove server with application on it")
	}
	// TODO: remove all connections fills for this server

	deployService.DatabaseAdapter.DeleteServer(server)
	return nil
}
