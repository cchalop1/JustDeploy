package application

import (
	"errors"
	"strings"

	"cchalop1.com/deploy/internal/adapter/database"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/internal/utils"
)

func CreateServiceFromDatabase(deployService *service.DeployService, databaseName string) (domain.Service, error) {
	// Get the database configuration
	dbConfig, err := database.GetServiceByName(databaseName)
	if err != nil {
		return domain.Service{}, errors.New("Database service not found")
	}

	// Generate a unique ID for the service
	serviceId := utils.GenerateRandomPassword(5)

	// Create a new service
	service := domain.Service{
		Id:         serviceId,
		Status:     "ready_to_deploy",
		Name:       strings.ToLower(databaseName),
		Type:       "database",
		ImageName:  dbConfig.Config.Image,
		ImageUrl:   dbConfig.Icon,
		IsRepo:     false,
		Envs:       dbConfig.Env,
		ExposePort: string(dbConfig.DefaultPort),
		ExposeSettings: domain.ServiceExposeSettings{
			IsExposed: false,
		},
	}

	// Save the service to the database
	deployService.DatabaseAdapter.SaveService(service)

	return service, nil
}
