package application

import (
	"cchalop1.com/deploy/internal/api/service"
)

func GithubIsConnected(deployService *service.DeployService) bool {
	settings := deployService.DatabaseAdapter.GetSettings()

	githubToken := settings.GithubToken

	return githubToken != ""
}
