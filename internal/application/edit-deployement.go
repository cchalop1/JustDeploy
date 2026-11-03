package application

import (
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
)

func EditDeploy(ds *service.DeployService, edit dto.EditDeployDto) error {
	return domain.ErrUnsupported
}
