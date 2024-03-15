package main

import (
	"os"

	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"cchalop1.com/deploy/internal/web"
)

func main() {
	app := api.NewApplication()

	databaseAdapter := adapter.NewDatabaseAdapter()
	filesystemAdapter := adapter.NewFilesystemAdapter()

	deployConfig := application.GetDeployConfig(databaseAdapter, filesystemAdapter)

	deployService := service.DeployService{
		DeployConfig:      &deployConfig,
		DatabaseAdapter:   databaseAdapter,
		DockerAdapter:     adapter.NewDockerAdapter(),
		FilesystemAdapter: filesystemAdapter,
	}

	// get the connection to the server if is exist
	//TODO: extract to a function
	if deployService.DeployConfig.DeployStatus != "serverconfig" {
		deployService.DockerAdapter = application.ConnectAndSetupServer(&deployService)
	}

	args := os.Args[1:]

	if len(args) > 0 {
		if args[0] == "redeploy" {
			application.ReDeployApplication(&deployService, deployConfig.AppConfig.Name)
			os.Exit(0)
		}
	} else {
		api.CreateRoutes(app, &deployService)
		web.CreateMiddlewareWebFiles(app)
		app.StartServer(true)
	}

}
