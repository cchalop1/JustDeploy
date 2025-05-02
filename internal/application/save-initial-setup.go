package application

import (
	"errors"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
)

// SaveInitialSetup saves the API key and domain during initial setup
func SaveInitialSetup(deployService *service.DeployService, setupDto dto.InitialSetupDto) error {
	// Get current settings
	settings := deployService.DatabaseAdapter.GetSettings()

	// Check if API key already exists
	if settings.ApiKey != "" {
		// If API key already exists, don't allow changing it
		return errors.New("API key already set and cannot be changed through this endpoint")
	}

	// First, save the API key to settings
	settings.ApiKey = setupDto.ApiKey

	err := deployService.DatabaseAdapter.SaveSettings(settings)
	if err != nil {
		return err
	}

	// Then, save the domain to the server
	server := deployService.DatabaseAdapter.GetServer()
	server.Domain = setupDto.Domain

	err = deployService.DatabaseAdapter.SaveServer(server)
	if err != nil {
		return err
	}

	return nil
}
