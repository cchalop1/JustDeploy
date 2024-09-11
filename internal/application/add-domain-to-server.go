package application

import (
	"fmt"

	"cchalop1.com/deploy/internal/api/http/dto"
	"cchalop1.com/deploy/internal/api/service"
)

func AddDomainToServer(deployService *service.DeployService, newDomain dto.NewDomain, serverId string) error {
	server, err := deployService.DatabaseAdapter.GetServerById(serverId)
	if err != nil {
		return err
	}

	server.Domain = newDomain.Domain
	fmt.Println(server)
	err = deployService.DatabaseAdapter.SaveServer(server)
	return err
}
