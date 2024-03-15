package application

import (
	"fmt"

	"cchalop1.com/deploy/internal/api/service"
)

func ReDeployApplication(deployService *service.DeployService, containerName string) {
	RemoveApplication(containerName, deployService)
	DeployApplication(deployService)
	fmt.Println("Success to redeploy ", containerName)
}
