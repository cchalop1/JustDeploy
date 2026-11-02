package application

import (
	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"fmt"
	"strconv"
)

func ConnectGithubApp(deployService *service.DeployService, code string) (adapter.GitHubAppResponse, error) {
	response, err := deployService.GithubAdapter.FinalizeGithubAppCreation(code)
	if err != nil {
		return adapter.GitHubAppResponse{}, err
	}

	err = saveGithubAppDetails(deployService, response)
	if err != nil {
		return adapter.GitHubAppResponse{}, err
	}

	return response, nil
}

func SaveAccessTokenWithInstallationId(deployService *service.DeployService, installationId string) error {
	if n, err := strconv.ParseInt(installationId, 10, 64); err != nil || n < 1 {
		return fmt.Errorf("invalid installation ID")
	}
	settings := deployService.DatabaseAdapter.GetSettings()
	appID := settings.GithubAppSettings.ID
	privateKey := settings.GithubAppSettings.Pem

	accessToken, err := deployService.GithubAdapter.GetInstallationToken(appID, privateKey, installationId)
	if err != nil {
		return err
	}

	settings.GithubToken = accessToken
	settings.GithubAppSettings.InstallationID = installationId
	err = deployService.DatabaseAdapter.SaveSettings(settings)
	if err != nil {
		return err
	}

	return nil
}

func saveGithubAppDetails(deployService *service.DeployService, response adapter.GitHubAppResponse) error {
	settings := deployService.DatabaseAdapter.GetSettings()

	settings.GithubAppSettings = domain.GithubAppSettings{
		ClientID:      response.ClientID,
		ClientSecret:  response.ClientSecret,
		WebhookSecret: response.WebhookSecret,
		Pem:           response.Pem,
		ID:            response.ID,
		Name:          response.Name,
	}

	err := deployService.DatabaseAdapter.SaveSettings(settings)
	if err != nil {
		return err
	}

	return nil
}
