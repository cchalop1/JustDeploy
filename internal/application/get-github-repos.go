package application

import (
	"fmt"

	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api/service"
)

func GetGithubRepos(deployService *service.DeployService) []adapter.GithubRepo {
	settings := deployService.DatabaseAdapter.GetSettings()

	githubToken := settings.GithubToken

	repos, err := deployService.GithubAdapter.GetRepos(githubToken)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return repos
}
