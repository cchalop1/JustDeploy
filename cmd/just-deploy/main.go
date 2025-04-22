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

// @title           JustDeploy API
// @version         1.0
// @description     JustDeploy is a PaaS tool designed to simplify the lives of developers. It allows you to easily deploy your projects and databases using Docker. JustDeploy fetches your GitHub repository and deploys your application using your Docker and Docker Compose configurations, all while deploying to any VPS of your choice without vendor lock-in.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Clément Chalopin
// @contact.url    https://github.com/cchalop1/JustDeploy
// @contact.email  support@justdeploy.dev

// @license.name  AGPL-3.0 License
// @license.url   https://github.com/cchalop1/JustDeploy/blob/main/LICENSE

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  GitHub Repository
// @externalDocs.url          https://github.com/cchalop1/JustDeploy
func main() {
	isNewVersion := application.CheckIsNewVersion()

	if isNewVersion {
		fmt.Println("There is a new version available. Please download it by typing: curl -fsSL https://raw.githubusercontent.com/cchalop1/JustDeploy/refs/heads/main/install.sh | bash")
	}

	port := adapter.FindOpenLocalPort(8080)

	app := api.NewApplication(port)

	api.CreateSwaggerRoutes(app)

	databaseAdapter := adapter.NewDatabaseAdapter()
	filesystemAdapter := adapter.NewFilesystemAdapter()
	dockerAdapter := adapter.NewDockerAdapter()
	networkAdapter := adapter.NewNetworkAdapter()
	githubAdapter := adapter.NewGithubAdapter()
	gitAdapter := adapter.NewGitAdapter()

	databaseAdapter.Init()

	deployService := service.DeployService{
		DatabaseAdapter:   databaseAdapter,
		DockerAdapter:     dockerAdapter,
		FilesystemAdapter: filesystemAdapter,
		EventAdapter:      adapter.NewAdapterEvent(),
		NetworkAdapter:    networkAdapter,
		GithubAdapter:     githubAdapter,
		GitAdapter:        gitAdapter,
	}

	getArgsOptions()

	server, apiKey, err := application.CreateCurrentServer(&deployService, port)
	isNewServer := err == nil

	if err != nil {
		fmt.Println("Current Server is arealy created :", err)
	}

	if flags.help {
		showHelp()
	}

	api.InitValidator(app)
	web.CreateMiddlewareWebFiles(app)
	api.CreateRoutes(app, &deployService)
	displayServerURL(networkAdapter, server, isNewServer, apiKey)
	app.StartServer(port)
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

func displayServerURL(networkAdapter *adapter.NetworkAdapter, server domain.Server, isNewServer bool, apiKey string) {
	url, err := networkAdapter.GetServerURL(server.Port)
	if err != nil {
		fmt.Println("Error getting server URL:", err)
		return
	}

	// Add welcome parameter if this is a new server
	if isNewServer {
		url = url + "?welcome=true"
	}

	// Add API key parameter if available
	if apiKey != "" {
		if isNewServer {
			url = url + "&api_key=" + apiKey
		} else {
			url = url + "?api_key=" + apiKey
		}
	}

	yellow := "\033[33m"
	reset := "\033[0m"

	fmt.Println(yellow + `
           _           _   _____             _             
          | |         | | |  __ \           | |            
          | |_   _ ___| |_| |  | | ___ _ __ | | ___  _   _ 
      _   | | | | / __| __| |  | |/ _ \ '_ \| |/ _ \| | | |
     | |__| | |_| \__ \ |_| |__| |  __/ |_) | | (_) | |_| |
      \____/ \__,_|___/\__|_____/ \___| .__/|_|\___/ \__, |
                                      | |             __/ |
                                      |_|            |___/ 
 	🛵 JustDeploy - Simplify Development & Deployment 🚀
	` + reset)

	fmt.Printf("Server is running at %s\n", url)
}
