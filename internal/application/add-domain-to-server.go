package application

import (
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
)

func AddDomainToServer(deployService *service.DeployService, newDomain dto.NewDomain) error {
	server := deployService.DatabaseAdapter.GetServer()
	server.Domain = newDomain.Domain
	err := deployService.DatabaseAdapter.SaveServer(server)
	return err
}
