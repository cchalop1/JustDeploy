package application

import "cchalop1.com/deploy/internal/api/service"

func SaveGithubToken(deployService *service.DeployService, token string) {
	settings := deployService.DatabaseAdapter.GetSettings()

	settings.GithubToken = token

	deployService.DatabaseAdapter.SaveSettings(settings)
}
