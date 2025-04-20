package application

import (
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/utils"
)

// GenerateAndSaveApiKey generates a new API key and saves it to the database settings
// It returns the generated API key
func GenerateAndSaveApiKey(deployService *service.DeployService) (string, error) {
	// Generate a random string for API key (32 characters)
	apiKey := utils.GenerateRandomPassword(32)

	// Get current settings
	settings := deployService.DatabaseAdapter.GetSettings()

	// Update settings with new API key
	settings.ApiKey = apiKey

	// Save updated settings
	err := deployService.DatabaseAdapter.SaveSettings(settings)
	if err != nil {
		return "", err
	}

	return apiKey, nil
}

// GetExistingApiKey returns the existing API key or empty string if not set
func GetExistingApiKey(deployService *service.DeployService) string {
	settings := deployService.DatabaseAdapter.GetSettings()
	return settings.ApiKey
}
