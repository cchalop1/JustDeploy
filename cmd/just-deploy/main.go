package main

import (
	"flag"
	"fmt"
	"os"

	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api"
	"cchalop1.com/deploy/internal/api/dto"
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
	app := api.NewApplication()

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

	// TODO: try server connection
	// TODO: do health check

	currentPath := filesystemAdapter.GetCurrentPath()

	project := databaseAdapter.GetProjectByPath(currentPath)

	if project == nil {
		createProjectDto := dto.CreateProjectDto{
			Name: filesystemAdapter.GetFolderName(currentPath),
			Path: currentPath,
		}

		projectId, err := application.CreateProject(&deployService, createProjectDto)
		if err != nil {
			fmt.Println("Erreur lors de la création du projet:", err)
			os.Exit(1)
		}

		fmt.Printf("Projet créé avec succès avec l'ID: %s\n", projectId)
	} else {
		fmt.Printf("Projet déjà existant: %s\n", project.Name)
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
		web.CreateMiddlewareWebFiles(app)
		if !flags.noBrowser {
			fmt.Println("Opening browser")
			adapter.OpenBrowser("http://localhost:8080/project/" + project.Id)
		}
		app.StartServer()

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
