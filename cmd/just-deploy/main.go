package main

import (
	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api"
	"cchalop1.com/deploy/internal/api/usecase"
	"cchalop1.com/deploy/internal/application"
	"cchalop1.com/deploy/internal/web"
)

func main() {
	app := api.NewApplication()

	databaseAdapter := adapter.NewDatabaseAdapter()
	deployConfig := application.GetDeployConfig(databaseAdapter)

	deployUseCase := usecase.DeployUseCase{
		DeployConfig:  &deployConfig,
		DockerAdapter: adapter.NewDockerAdapter(),
	}

	// get the connection to the server if is exist
	//TODO: extract to a function
	if deployUseCase.DeployConfig.DeployStatus != "serverconfig" {
		deployUseCase.DockerAdapter = application.ConnectAndSetupServer(deployUseCase.DeployConfig.ServerConfig)
	}

	api.CreateRoutes(app, &deployUseCase)
	web.CreateMiddlewareWebFiles(app)
	app.StartServer(true)
}
