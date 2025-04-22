package application

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"cchalop1.com/deploy/internal/adapter/database"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/internal/utils"
)

// CreateServiceFromDatabase creates a new database service based on preconfigured templates
// from the database adapter.
func CreateServiceFromDatabase(deployService *service.DeployService, databaseName string) (domain.Service, error) {
	// Get the database configuration
	dbConfig, err := database.GetServiceByName(databaseName)
	if err != nil {
		return domain.Service{}, errors.New("Database service not found")
	}

	// Generate a unique ID for the service
	serviceId := utils.GenerateRandomPassword(5)

	// Generate values for environment variables
	envsWithValues := utils.GenerateEnvValues(dbConfig.Env)

	// Create a new service for the database
	service := domain.Service{
		Id:           serviceId,
		Status:       "ready_to_deploy",
		Name:         strings.ToLower(databaseName) + "-" + serviceId,
		Type:         "database",
		ImageName:    dbConfig.Config.Image,
		ImageUrl:     dbConfig.Icon,
		Envs:         envsWithValues,
		DockerHubUrl: utils.GetDockerHubUrl(dbConfig.Config.Image),
		ExposeSettings: domain.ServiceExposeSettings{
			IsExposed:  false,
			ExposePort: strconv.Itoa(dbConfig.DefaultPort),
		},
	}

	// Log the service creation
	fmt.Printf("Creating new database service: %s (type: %s, image: %s)\n",
		service.Name, service.Type, service.ImageName)

	// Save the service to the database
	deployService.DatabaseAdapter.SaveService(service)

	return service, nil
}
