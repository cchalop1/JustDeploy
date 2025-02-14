package application

import (
	"cchalop1.com/deploy/internal/adapter/database"
	"cchalop1.com/deploy/internal/api/service"
)

func GetConfiguredServiceList(deployService *service.DeployService, productId string) []database.ServicesConfig {
	services := database.GetListOfDatabasesServices()

	if productId == "" {
		return services
	}

	project, err := deployService.DatabaseAdapter.GetProjectById(productId)

	if err != nil {
		return services
	}

	servicesList := []database.ServicesConfig{}

	for _, service := range services {
		keepService := true
		for _, projectService := range project.Services {
			if service.Config.Image == projectService.ImageName {
				keepService = false
			}
		}
		if keepService {
			servicesList = append(servicesList, service)
		}
	}
	return servicesList
}
