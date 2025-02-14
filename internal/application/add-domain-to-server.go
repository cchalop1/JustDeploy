package application

import (
	"fmt"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
)

func AddDomainToServer(deployService *service.DeployService, newDomain dto.NewDomain, serverId string) error {
	server := deployService.DatabaseAdapter.GetServer()

	server.Domain = newDomain.Domain
	fmt.Println(server)
	err := deployService.DatabaseAdapter.SaveServer(server)
	return err
}
