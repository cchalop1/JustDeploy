package application

import (
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
)

func GetServicesFromDockerCompose(deployService *service.DeployService, deployId string) ([]dto.ServiceDto, error) {
	deploy, err := deployService.DatabaseAdapter.GetDeployById(deployId)
	if err != nil {
		return []dto.ServiceDto{}, err
	}

	services, err := deployService.FilesystemAdapter.GetComposeConfigOfDeploy(deploy.PathToSource)

	return services, nil
}
