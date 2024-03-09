// internal/main.go
package main

import (
	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/application"
)

func main() {
	filesystemAdapter := adapter.NewFilesystemAdapter()

	deployConfig := application.GetFormDetails(filesystemAdapter)
	httpAdapter := application.NewHttpAdapter(deployConfig)

	httpAdapter.StartServer(true)

	// dockerAdapter := adapter.NewDockerAdapter()

	// sshAdapter := adapter.NewSshAdapter()

	// setupVmService := application.NewSetupVMService(sshAdapter)

	// setupVmService.Setup(deployConfig)

	// deployementService := application.NewDeploymentService(dockerAdapter, sshAdapter)

	// deployementService.DeployApplication(deployConfig)
}
