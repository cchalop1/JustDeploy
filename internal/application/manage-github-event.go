package application

import (
	"fmt"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
)

func ManageGithubEvent(deployService *service.DeployService, githubEvent dto.GithubEvent) error {
	var err error
	fmt.Println(githubEvent)

	if githubEvent.Ref == "refs/heads/main" {
		err = ReDeployApplication(deployService, githubEvent.Repository.Name)
	}
	return err
}
