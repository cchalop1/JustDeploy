package application

import (
	"fmt"

	"cchalop1.com/deploy/internal/api/service"
)

func ReDeployApplication(deployService *service.DeployService, containerName string) {
	// TODO: implement this without deplouConfig
	RemoveApplication(containerName, deployService)
	// DeployApplication(deployService)
	fmt.Println("Success to redeploy ", containerName)
}
