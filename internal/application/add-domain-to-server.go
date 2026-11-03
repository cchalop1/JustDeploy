package application

import (
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"errors"
)

func AddDomainToServer(deployService *service.DeployService, newDomain dto.NewDomain) error {
	if !ValidDomain(newDomain.Domain) {
		return errors.New("invalid domain")
	}
	server := deployService.DatabaseAdapter.GetServer()
	server.Domain = newDomain.Domain
	err := deployService.DatabaseAdapter.SaveServer(server)
	return err
}
