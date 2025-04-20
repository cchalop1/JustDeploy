package application

import (
	"time"

	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/internal/utils"
)

func CreateCurrentServer(deployService *service.DeployService, port string) (domain.Server, string, error) {
	oldServer := deployService.DatabaseAdapter.GetServer()
	currentIp, err := deployService.NetworkAdapter.GetCurrentIP()

	if oldServer.Ip == currentIp {
		// Existing server, check for API key
		settings := deployService.DatabaseAdapter.GetSettings()
		if settings.ApiKey == "" {
			// Generate API key for existing server
			apiKey, err := GenerateAndSaveApiKey(deployService)
			if err != nil {
				return oldServer, "", err
			}
			return oldServer, apiKey, nil
		}
		return oldServer, settings.ApiKey, nil
	}

	if err != nil {
		return domain.Server{}, "", err
	}

	server := domain.Server{
		Id:          utils.GenerateRandomPassword(5),
		Name:        "Local Server",
		Ip:          currentIp,
		Port:        port,
		CreatedDate: time.Now(),
		Status:      "Installing",
	}

	deployService.DatabaseAdapter.SaveServer(server)

	// Generate API key for new server
	apiKey, err := GenerateAndSaveApiKey(deployService)
	if err != nil {
		return server, "", err
	}

	return server, apiKey, nil
}
