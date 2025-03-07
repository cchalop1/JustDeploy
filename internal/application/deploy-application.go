package application

import (
	"errors"
	"fmt"
	"path/filepath"

	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
)

func deployGithubService(deployService *service.DeployService, serviceToDeploy domain.Service, baseDomain string) error {
	fmt.Println("service to deploy : ", serviceToDeploy, baseDomain)
	pathToDir, err := filepath.Abs(serviceToDeploy.CurrentPath)

	if err != nil {
		return err
	}

	pathToDir = adapter.NewFilesystemAdapter().CleanPath(pathToDir)

	isFolder := adapter.NewFilesystemAdapter().IsFolder(pathToDir)
	DockerFileName := "Dockerfile"

	if !isFolder {
		fmt.Println("Is not a folder")
		DockerFileName = adapter.NewFilesystemAdapter().BaseDir(pathToDir)
		pathToDir = adapter.NewFilesystemAdapter().GetDir(pathToDir)
	}

	portEnv := make([]dto.Env, 1)

	portEnv[0] = dto.Env{
		Name:  "PORT",
		Value: "80",
	}

	fmt.Println("Path to dir: ", pathToDir)
	fmt.Println("Docker file name: ", DockerFileName)

	serviceToDeploy.Envs = append(portEnv, serviceToDeploy.Envs...)

	isDockerfile := deployService.FilesystemAdapter.FindDockerFile(serviceToDeploy.CurrentPath)
	if isDockerfile {
		err := deployService.DockerAdapter.BuildImage(serviceToDeploy)
		if err != nil {
			return fmt.Errorf("error building Docker image: %w", err)
		}
	} else {
		err = deployService.DockerAdapter.BuildNixpacksImage(serviceToDeploy)
		if err != nil {
			return fmt.Errorf("error building Nixpacks image: %w", err)
		}
	}

	serviceDomain := serviceToDeploy.Name + "." + baseDomain

	err = deployService.DockerAdapter.RunImage(serviceToDeploy, serviceDomain)

	if err != nil {
		return fmt.Errorf("error running Docker image: %w", err)
	}

	serviceToDeploy.Status = "Running"
	serviceToDeploy.SetUrl(baseDomain)
	deployService.DatabaseAdapter.SaveService(serviceToDeploy)

	return nil
}

func deployDatabaseService(deployService *service.DeployService, dbService domain.Service) error {
	isRunning, err := deployService.DockerAdapter.IsServiceRunning(dbService.GetDockerName())
	if err != nil {
		return fmt.Errorf("error checking if database service %s is running: %w", dbService.GetDockerName(), err)
	}

	if isRunning {
		fmt.Printf("Database service %s is already running\n", dbService.GetDockerName())
		return nil
	}

	fmt.Printf("Deploying database service: %s\n", dbService.GetDockerName())

	// Deploy the database service
	deployService.DockerAdapter.RunServiceWithDeploy(dbService, dbService.GetDockerName())

	dbService.Status = "Running"

	// Save the service to the database
	deployService.DatabaseAdapter.SaveService(dbService)

	return nil
}

func deployOneService(deployService *service.DeployService, serviceToDeploy domain.Service, baseDomain string) error {
	fmt.Println(serviceToDeploy)
	switch serviceToDeploy.Type {
	case "database":
		return deployDatabaseService(deployService, serviceToDeploy)
	case "github_repo":
		return deployGithubService(deployService, serviceToDeploy, baseDomain)
	default:
		// For backward compatibility or other service types, default to GitHub repo deployment
		return deployGithubService(deployService, serviceToDeploy, baseDomain)
	}
}

func DeployApplication(deployService *service.DeployService) error {
	services := deployService.DatabaseAdapter.GetServices()

	server := deployService.DatabaseAdapter.GetServer()

	err := deployService.DockerAdapter.ConnectClient(server)

	if err != nil {
		return fmt.Errorf("error connecting Docker client: %w", err)
	}

	// Check if Traefik router is running, if not, pull and run it
	routerRunning, err := deployService.DockerAdapter.CheckRouterIsRunning()
	if err != nil {
		return fmt.Errorf("error checking if router is running: %w", err)
	}

	if !routerRunning {
		err = deployService.DockerAdapter.PullTreafikImage()
		if err != nil {
			return fmt.Errorf("error pulling Traefik image: %w", err)
		}

		err = deployService.DockerAdapter.RunRouter("clement.chalopin@gmail.com")
		if err != nil {
			return fmt.Errorf("error running Traefik router: %w", err)
		}
	} else {
		fmt.Println("Traefik router is already running")
	}

	if server.Domain == "" {
		return errors.New("server does not have domain")
	}

	// Deploy application services that are not already running
	for _, service := range services {
		isRunning, err := deployService.DockerAdapter.IsServiceRunning(service.GetDockerName())
		if err != nil {
			return fmt.Errorf("error checking if service %s is running: %w", service.Name, err)
		}

		if !isRunning {
			fmt.Printf("Deploying service: %s\n", service.Name)
			err = deployOneService(deployService, service, server.Domain)
			if err != nil {
				return fmt.Errorf("error deploying service: %w", err)
			}
		} else {
			fmt.Printf("Service %s is already running\n", service.Name)
		}
	}

	// Deploy all database services from the database configuration
	// databaseServices := database.GetListOfDatabasesServices()
	// for _, dbService := range databaseServices {
	// 	err = deployDatabaseService(deployService, dbService.Name)
	// 	if err != nil {
	// 		return fmt.Errorf("error deploying database service: %w", err)
	// 	}
	// }

	return nil
}
