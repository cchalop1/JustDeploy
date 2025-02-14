package application

import (
	"cchalop1.com/deploy/internal/adapter/database"
	"cchalop1.com/deploy/internal/api/service"
)

func GetConfiguredServiceList(deployService *service.DeployService, productId string) []database.ServicesConfig {
	services := database.GetListOfDatabasesServices()

	return services
}
