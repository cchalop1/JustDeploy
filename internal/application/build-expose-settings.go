package application

import (
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
)

// BuildExposeSettings creates the appropriate ExposeSettings for a service
// If it's the first github_repo service, the subdomain will be empty
// Otherwise, the subdomain will be set to the service's Name
func BuildExposeSettings(deployService *service.DeployService, serviceName string, isExposed bool, exposePort string) domain.ServiceExposeSettings {
	// Get all existing services
	existingServices := deployService.DatabaseAdapter.GetServices()

	// Check if there are any existing github_repo services
	isFirstGithubRepo := true
	for _, svc := range existingServices {
		if svc.Type == "github_repo" {
			isFirstGithubRepo = false
			break
		}
	}

	// Create the ExposeSettings
	exposeSettings := domain.ServiceExposeSettings{
		IsExposed:  isExposed,
		ExposePort: exposePort,
		Tls:        false, // Default value, can be made configurable if needed
	}

	// Set the SubDomain based on whether it's the first github_repo service
	if isFirstGithubRepo {
		exposeSettings.SubDomain = "" // Empty subdomain for the first github_repo
	} else {
		exposeSettings.SubDomain = serviceName // Use service name for subsequent services
	}

	return exposeSettings
}
