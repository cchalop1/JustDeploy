package application

import (
	"cchalop1.com/deploy/internal/adapter/database"
	"cchalop1.com/deploy/internal/api/service"
)

func GetServiceList(deployService *service.DeployService) []database.ServicesConfig {
	services := database.GetListOfDatabasesServices()
	return services
}
