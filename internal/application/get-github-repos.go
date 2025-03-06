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
	// Get installation ID
	installationID, err := deployService.GithubAdapter.GetInstallationID(appSettings.ID, appSettings.Pem)
	if err != nil {
		return "", fmt.Errorf("failed to get installation ID: %w", err)
	}

	// Convert installation ID to string
	installationIDStr := fmt.Sprintf("%d", installationID)

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

func GetGithubRepos(deployService *service.DeployService) []adapter.GithubRepo {
	settings := deployService.DatabaseAdapter.GetSettings()
	githubToken := settings.GithubToken
	appSettings := settings.GithubAppSettings

	// Check if token is valid before making the API call
	if !deployService.GithubAdapter.IsTokenValid(githubToken) {
		// Token is invalid or expired, renew it
		newToken, err := renewGithubToken(deployService, appSettings)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		githubToken = newToken
	}

	// Get installation ID for token renewal if needed
	installationID, err := deployService.GithubAdapter.GetInstallationID(appSettings.ID, appSettings.Pem)
	if err != nil {
		fmt.Println("Failed to get installation ID:", err)
		return nil
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
		fmt.Println(err)
		return nil
	}

	return repos
}
