package main

import (
	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"cchalop1.com/deploy/internal/web"
)

func main() {
	app := api.NewApplication()

	databaseAdapter := adapter.NewDatabaseAdapter()
	deployConfig := application.GetDeployConfig(databaseAdapter)

	deployService := service.DeployService{
		DeployConfig:    &deployConfig,
		DatabaseAdapter: databaseAdapter,
		DockerAdapter:   adapter.NewDockerAdapter(),
	}

	// get the connection to the server if is exist
	//TODO: extract to a function
	if deployService.DeployConfig.DeployStatus != "serverconfig" {

		deployService.DockerAdapter = application.ConnectAndSetupServer(&deployService)
	}

	api.CreateRoutes(app, &deployService)
	web.CreateMiddlewareWebFiles(app)
	app.StartServer(true)
}
