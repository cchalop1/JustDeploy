package application

import (
	"cchalop1.com/deploy/internal/adapter/database"
	"cchalop1.com/deploy/internal/api/service"
)

func GetServicesFromDockerCompose(deployService *service.DeployService, deployId string) ([]database.ServicesConfig, error) {
	// deploy, err := deployService.DatabaseAdapter.GetDeployById(deployId)
	// if err != nil {
	// 	return []database.ServicesConfig{}, err
	// }

	// services, err := deployService.FilesystemAdapter.GetComposeConfigOfDeploy(deploy.PathToSource)

	// if err != nil {
	// 	return []database.ServicesConfig{}, err
	// }

	// return services, nil
	return nil, nil
}
