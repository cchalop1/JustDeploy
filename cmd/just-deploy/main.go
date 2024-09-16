package main

import (
	"flag"
	"fmt"
	"os"

	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"cchalop1.com/deploy/internal/web"
)

var flags struct {
	noBrowser bool
	help      bool
	redeploy  struct {
		deployId string
	}
}

func main() {
	currentVersion := application.GetVersion()
	fmt.Println("JustDeploy version: ", currentVersion.Version)

	port := adapter.FindOpenLocalPort(8080)

	app := api.NewApplication(port)

	databaseAdapter := adapter.NewDatabaseAdapter()
	filesystemAdapter := adapter.NewFilesystemAdapter()
	dockerAdapter := adapter.NewDockerAdapter()

	databaseAdapter.Init()

	deployService := service.DeployService{
		DatabaseAdapter:   databaseAdapter,
		DockerAdapter:     dockerAdapter,
		FilesystemAdapter: filesystemAdapter,
		EventAdapter:      adapter.NewAdapterEvent(),
	}

	projectId, err := application.CreateProjectCurrentFolder(&deployService)

	if err != nil {
		fmt.Println("Error while creating project:", err)
		os.Exit(1)
	}

	getArgsOptions()

	if flags.help {
		showHelp()
	}

	if flags.redeploy.deployId != "" {
		application.ReDeployApplication(&deployService, flags.redeploy.deployId)
		os.Exit(0)
	} else {
		api.InitValidator(app)
		api.CreateRoutes(app, &deployService)
		// web.ReplaceEnvInEnvBuild(port)
		web.CreateMiddlewareWebFiles(app)
		if !flags.noBrowser {
			fmt.Println("Opening browser")
			adapter.OpenBrowser("http://localhost:" + port + "/project/" + projectId)
		}
		app.StartServer(port)
	}
}

func getArgsOptions() {
	flag.BoolVar(&flags.noBrowser, "no-browser", false, "Do not open the browser")
	flag.StringVar(&flags.redeploy.deployId, "redeploy", "", "Redeploy application by deploy id")
	flag.BoolVar(&flags.help, "help", false, "Show help")
	flag.Parse()
}

func showHelp() {
	fmt.Println("Usage: main [options]")
	fmt.Println("  -no-browser    Do not open the browser")
	fmt.Println("  -redeploy <id> Redeploy application by deploy id")
	os.Exit(0)
}
