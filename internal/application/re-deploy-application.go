package application

import (
	"errors"
	"fmt"
	"path/filepath"

	"cchalop1.com/deploy/internal/api/service"
)

func ReDeployApplication(deployService *service.DeployService, serviceName string) error {
	var err error

	settings := deployService.DatabaseAdapter.GetSettings()
	if settings.GithubToken == "" {
		return errors.New("GitHub token not found. Please configure GitHub integration in settings")
	}

	// Check if token is valid, renew it if necessary
	githubToken := settings.GithubToken
	appSettings := settings.GithubAppSettings

	// Check if token is valid before using it
	if !deployService.GithubAdapter.IsTokenValid(githubToken) {
		// Token is invalid or expired, renew it
		newToken, err := renewGithubToken(deployService, appSettings)
		if err != nil {
			return errors.New("Failed to renew GitHub token: " + err.Error())
		}
		githubToken = newToken

		// Update the settings with the new token
		settings.GithubToken = newToken
		err = deployService.DatabaseAdapter.SaveSettings(settings)
		if err != nil {
			return errors.New("Failed to save new GitHub token: " + err.Error())
		}
	}

	services := GetServices(deployService)

	for _, s := range services {
		if s.Name == serviceName {
			tempDir := deployService.FilesystemAdapter.GetTempDir()
			repoPath := filepath.Join(tempDir, s.Name)

			s.CurrentPath = repoPath

			err = deployService.GitAdapter.CloneRepository(s.FullName, s.CurrentPath, githubToken)
			if err != nil {
				return err
			}

			// Mettre à jour les informations du commit après le clonage
			err = UpdateServiceCommitInfo(deployService, &s)
			if err != nil {
				// Log l'erreur mais ne pas faire échouer le redéploiement
				fmt.Printf("Warning: failed to update commit info for service %s: %v\n", s.Name, err)
			}

			err = deployService.DatabaseAdapter.SaveService(s)
			if err != nil {
				return err
			}
			break
		}
	}

	err = DeployApplication(deployService)
	return err
}
