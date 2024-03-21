package application

import (
	"fmt"

	"cchalop1.com/deploy/internal/api/service"
)

func ReDeployApplication(deployService *service.DeployService, deployId string) error {
	deploy, err := deployService.DatabaseAdapter.GetDeployById(deployId)
	if err != nil {
		return err
	}
	RemoveApplicationById(deployService, deploy.Name)
	err = ReDeployApplicationRun(deployService, deploy)

	if err != nil {
		return err
	}

	fmt.Println("Success to redeploy ", deploy.Name)
	return nil
}
