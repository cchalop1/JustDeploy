package application

import (
	"fmt"
	"strconv"

	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
)

// renewGithubToken handles the GitHub token renewal process
func renewGithubToken(deployService *service.DeployService, appSettings domain.GithubAppSettings) (string, error) {
	installationIDStr := appSettings.InstallationID
	if installationIDStr == "" {
		installationID, err := deployService.GithubAdapter.GetInstallationID(appSettings.ID, appSettings.Pem)
		if err != nil {
			return "", err
		}
		installationIDStr = strconv.Itoa(installationID)
	}

	// Renew token
	newToken, err := deployService.GithubAdapter.RenewGithubToken(appSettings.ID, appSettings.Pem, installationIDStr)
	if err != nil {
		return "", fmt.Errorf("failed to renew GitHub token: %w", err)
	}

	// Save the new token
	err = deployService.DatabaseAdapter.SaveInstallationToken(installationIDStr, newToken)
	if err != nil {
		return "", fmt.Errorf("failed to save new token: %w", err)
	}

	return newToken, nil
}

func GetGithubRepos(deployService *service.DeployService) ([]adapter.GithubRepo, error) {
	settings := deployService.DatabaseAdapter.GetSettings()
	githubToken := settings.GithubToken
	appSettings := settings.GithubAppSettings

	// Check if token is valid before making the API call
	if !deployService.GithubAdapter.IsTokenValid(githubToken) {
		// Token is invalid or expired, renew it
		newToken, err := renewGithubToken(deployService, appSettings)
		if err != nil {
			return nil, err
		}
		githubToken = newToken
	}

	// Get installation ID for token renewal if needed
	installationID, err := deployService.GithubAdapter.GetInstallationID(appSettings.ID, appSettings.Pem)
	if err != nil {
		return nil, err
	}
	installationIDStr := strconv.Itoa(installationID)

	// Use the convenience method that handles token renewal automatically
	// This will return only the 15 most recently updated repositories
	repos, err := deployService.GithubAdapter.GetReposWithTokenRenewal(
		githubToken,
		appSettings.ID,
		appSettings.Pem,
		installationIDStr,
		deployService.DatabaseAdapter.SaveInstallationToken,
	)

	if err != nil {
		return nil, err
	}

	return repos, nil
}
