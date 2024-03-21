package main

import (
	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/web"
)

func main() {
	app := api.NewApplication()

	databaseAdapter := adapter.NewDatabaseAdapter()
	filesystemAdapter := adapter.NewFilesystemAdapter()
	dockerAdapter := adapter.NewDockerAdapter()

	databaseAdapter.Init()

	deployService := service.DeployService{
		DatabaseAdapter:   databaseAdapter,
		DockerAdapter:     dockerAdapter,
		FilesystemAdapter: filesystemAdapter,
	}

	// TODO: try server connection
	// TODO: do health check

	// args := os.Args[1:]

	// if len(args) > 0 {
	// 	if args[0] == "redeploy" {
	// 		application.ReDeployApplication(&deployService, deployConfig.AppConfig.Name)
	// 		os.Exit(0)
	// 	}
	// } else {
	api.CreateRoutes(app, &deployService)
	web.CreateMiddlewareWebFiles(app)
	app.StartServer(false)
	// }

}
