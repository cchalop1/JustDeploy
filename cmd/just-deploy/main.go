package main

import (
	"flag"
	"fmt"
	"os"

	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"cchalop1.com/deploy/internal/domain"
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
	isNewVersion := application.CheckIsNewVersion()

	if isNewVersion {
		fmt.Println("There is a new version available. Please download it by typing: curl -fsSL https://get.justdeploy.app | bash")
	}

	port := adapter.FindOpenLocalPort(8080)

	app := api.NewApplication(port)

	databaseAdapter := adapter.NewDatabaseAdapter()
	filesystemAdapter := adapter.NewFilesystemAdapter()
	dockerAdapter := adapter.NewDockerAdapter()
	networkAdapter := adapter.NewNetworkAdapter()

	databaseAdapter.Init()

	deployService := service.DeployService{
		DatabaseAdapter:   databaseAdapter,
		DockerAdapter:     dockerAdapter,
		FilesystemAdapter: filesystemAdapter,
		EventAdapter:      adapter.NewAdapterEvent(),
		NetworkAdapter:    networkAdapter,
	}

	getArgsOptions()

	server, err := application.CreateCurrentServer(&deployService, port)

	if err != nil {
		fmt.Println("Current Server is arealy created :", err)
	}

	if flags.help {
		showHelp()
	}

	if flags.redeploy.deployId != "" {
		application.ReDeployApplication(&deployService, flags.redeploy.deployId)
		os.Exit(0)
	} else {
		api.InitValidator(app)
		api.CreateRoutes(app, &deployService)
		web.CreateMiddlewareWebFiles(app)
		displayServerURL(networkAdapter, server)
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

func displayServerURL(networkAdapter *adapter.NetworkAdapter, server domain.Server) {
	url, err := networkAdapter.GetServerURL(server.Port)
	if err != nil {
		fmt.Println("Error getting server URL:", err)
		return
	}
	fmt.Printf("Server is running at %s\n", url)
}
