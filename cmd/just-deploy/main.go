package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"cchalop1.com/deploy/internal"
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
	logs      bool
	update    bool
	stop      bool
	uninstall bool
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

// @host      localhost:5915
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  GitHub Repository
// @externalDocs.url          https://github.com/cchalop1/JustDeploy
func main() {
	isNewVersion := application.CheckIsNewVersion()

	if isNewVersion {
		fmt.Println("There is a new version available. Please download it by typing: curl -fsSL https://raw.githubusercontent.com/cchalop1/JustDeploy/refs/heads/main/install.sh | bash")
	}

	getArgsOptions()

	if flags.help {
		showHelp()
	}

	if flags.logs {
		showLogs()
		return
	}

	if flags.update {
		updateJustDeploy()
		return
	}

	if flags.stop {
		stopJustDeploy()
		return
	}

	if flags.uninstall {
		uninstallJustDeploy()
		return
	}

	port := "5915"

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

	deployService.DockerAdapter.ConnectClient()

	server, err := application.CreateCurrentServer(&deployService, port)

	if err != nil {
		fmt.Println("Current Server is already created:", err)
	}

	api.InitValidator(app)
	web.CreateMiddlewareWebFiles(app)
	api.CreateRoutes(app, &deployService)
	displayServerURL(networkAdapter, server)
	app.StartServer(port)
}

func getArgsOptions() {
	flag.BoolVar(&flags.noBrowser, "no-browser", false, "Do not open the browser")
	flag.StringVar(&flags.redeploy.deployId, "redeploy", "", "Redeploy application by deploy id")
	flag.BoolVar(&flags.help, "help", false, "Show help")
	flag.BoolVar(&flags.logs, "logs", false, "Show application logs")
	flag.BoolVar(&flags.update, "update", false, "Update JustDeploy to the latest version")
	flag.BoolVar(&flags.stop, "stop", false, "Stop JustDeploy service")
	flag.BoolVar(&flags.uninstall, "uninstall", false, "Uninstall JustDeploy completely")
	flag.Parse()
}

func showHelp() {
	fmt.Println("Usage: justdeploy [options]")
	fmt.Println("Options:")
	fmt.Println("  -no-browser    Do not open the browser")
	fmt.Println("  -redeploy <id> Redeploy application by deploy id")
	fmt.Println("  -logs          Show application logs")
	fmt.Println("  -update        Update JustDeploy to the latest version")
	fmt.Println("  -stop          Stop JustDeploy service")
	fmt.Println("  -uninstall     Uninstall JustDeploy completely")
	fmt.Println("  -help          Show this help message")
	os.Exit(0)
}

func showLogs() {
	cmd := exec.Command("sudo", "journalctl", "-u", "justdeploy.service", "-f")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("Showing JustDeploy logs (press Ctrl+C to exit)...")
	cmd.Run()
}

func updateJustDeploy() {
	fmt.Println("Updating JustDeploy...")
	cmd := exec.Command("bash", "-c", "curl -fsSL https://get.justdeploy.app | sudo bash")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error updating JustDeploy: %v\n", err)
		os.Exit(1)
	}
}

func stopJustDeploy() {
	fmt.Println("Stopping JustDeploy service...")
	cmd := exec.Command("sudo", "systemctl", "stop", "justdeploy.service")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error stopping JustDeploy service: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("JustDeploy service stopped successfully.")
}

func uninstallJustDeploy() {
	fmt.Println("Uninstalling JustDeploy...")

	// Stop and disable service
	stopCmd := exec.Command("sudo", "systemctl", "stop", "justdeploy.service")
	stopCmd.Run()

	disableCmd := exec.Command("sudo", "systemctl", "disable", "justdeploy.service")
	disableCmd.Run()

	// Remove service file
	rmServiceCmd := exec.Command("sudo", "rm", "-f", "/etc/systemd/system/justdeploy.service")
	rmServiceCmd.Run()

	// Reload systemd
	reloadCmd := exec.Command("sudo", "systemctl", "daemon-reload")
	reloadCmd.Run()

	// Remove binary
	rmBinCmd := exec.Command("sudo", "rm", "-f", "/usr/local/bin/justdeploy")
	rmBinCmd.Run()

	// Remove justdeploy folder and all its contents
	fmt.Println("Removing JustDeploy data directory...")
	err := os.RemoveAll(internal.JUSTDEPLOY_FOLDER)
	if err != nil {
		fmt.Printf("Error removing JustDeploy data directory: %v\n", err)
	}

	// Remove database file
	err = os.Remove(internal.DATABASE_FILE_PATH)
	if err != nil && !os.IsNotExist(err) {
		fmt.Printf("Error removing database file: %v\n", err)
	}

	// Remove cert docker folder
	err = os.RemoveAll(internal.CERT_DOCKER_FOLDER)
	if err != nil && !os.IsNotExist(err) {
		fmt.Printf("Error removing certificate folder: %v\n", err)
	}

	fmt.Println("JustDeploy has been completely uninstalled.")
}

func displayServerURL(networkAdapter *adapter.NetworkAdapter, server domain.Server) {
	url, err := networkAdapter.GetServerURL(server.Port)
	if err != nil {
		fmt.Println("Error getting server URL:", err)
		return
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
