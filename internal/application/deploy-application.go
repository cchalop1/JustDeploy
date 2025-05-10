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

// buildServiceDomain constructs the appropriate domain for a service based on its ExposeSettings.
//
// The domain is constructed as follows:
// - If ExposeSettings.SubDomain is empty: domain will be just the baseDomain (e.g. "example.com")
// - If ExposeSettings.SubDomain is provided: domain will be subdomain.baseDomain (e.g. "api.example.com")
//
// This allows services to be deployed either at the root domain or at a subdomain, consistent
// with how Service.SetUrl handles domain construction.
func buildServiceDomain(service domain.Service, baseDomain string) string {
	subDomain := service.ExposeSettings.SubDomain

	if subDomain == "" {
		return baseDomain // No subdomain - deploy at root domain (e.g. "example.com")
	}

	return subDomain + "." + baseDomain // With subdomain (e.g. "api.example.com")
}

func deployGithubService(deployService *service.DeployService, serviceToDeploy domain.Service, baseDomain string) error {
	// Stop any existing instance of the service
	deployService.DockerAdapter.Stop(serviceToDeploy.GetDockerName())

	serviceToDeploy.Status = "Deploying"
	deployService.DatabaseAdapter.SaveService(serviceToDeploy)

	fmt.Println("service to deploy : ", serviceToDeploy, baseDomain)
	pathToDir, err := filepath.Abs(serviceToDeploy.CurrentPath)

	if err != nil {
		return err
	}

	pathToDir = adapter.NewFilesystemAdapter().CleanPath(pathToDir)

	isFolder := adapter.NewFilesystemAdapter().IsFolder(pathToDir)

	if !isFolder {
		fmt.Println("Is not a folder")
		pathToDir = adapter.NewFilesystemAdapter().GetDir(pathToDir)
	}

	portEnv := make([]dto.Env, 1)
	serviceToDeploy.Envs = append(portEnv, serviceToDeploy.Envs...)

	logs := domain.LogCollector{}

	// Build the appropriate image based on whether a Dockerfile exists
	isDockerfile := deployService.FilesystemAdapter.FindDockerFile(serviceToDeploy.CurrentPath)
	if isDockerfile {
		err = deployService.DockerAdapter.BuildImage(serviceToDeploy, &logs)
	} else {
		err = deployService.DockerAdapter.BuildNixpacksImage(serviceToDeploy, &logs)
	}

	fmt.Println("logs : ", logs)

	fmt.Println("serviceToDeploy : ", serviceToDeploy)
	fmt.Println("logs : ", logs)

	StoreBuildLog(serviceToDeploy.Id, logs)

	if err != nil {
		fmt.Println("error building image: ", err)
		return fmt.Errorf("error building image: %w", err)
	}

	// Construct the domain for the service based on ExposeSettings
	serviceDomain := buildServiceDomain(serviceToDeploy, baseDomain)

	// Get the server to check HTTPS settings
	server := deployService.DatabaseAdapter.GetServer()

	// Run the container with the appropriate domain and TLS settings
	err = deployService.DockerAdapter.RunImageWithTLS(serviceToDeploy, serviceDomain, server.UseHttps)

	if err != nil {
		return fmt.Errorf("error running Docker image: %w", err)
	}

	// Update and save service status
	serviceToDeploy.Status = "Running"

	// Set the service URL - use baseDomain as SetUrl will handle subdomain logic internally
	serviceToDeploy.SetUrl(baseDomain)

	deployService.DatabaseAdapter.SaveService(serviceToDeploy)

	return nil
}

func deployDatabaseService(deployService *service.DeployService, dbService domain.Service) error {
	// Stop any existing instance of the service
	deployService.DockerAdapter.Stop(dbService.GetDockerName())

	fmt.Printf("Deploying database service: %s\n", dbService.GetDockerName())

	// Get the server to get HTTPS settings
	server := deployService.DatabaseAdapter.GetServer()

	// Deploy the database service (we pass empty domain since databases don't need external exposure)
	err := deployService.DockerAdapter.RunImageWithTLS(dbService, "", server.UseHttps)

	if err != nil {
		dbService.Status = "Error"
		deployService.DatabaseAdapter.SaveService(dbService)
		return fmt.Errorf("error running database container: %w", err)
	}

	// Update and save service status
	dbService.Status = "Running"
	deployService.DatabaseAdapter.SaveService(dbService)

	return nil
}

// deployOneService routes the service deployment to the appropriate handler based on service type
func deployOneService(deployService *service.DeployService, serviceToDeploy domain.Service, baseDomain string) error {
	fmt.Printf("Initiating deployment of service: %s (type: %s)\n", serviceToDeploy.Name, serviceToDeploy.Type)

	switch serviceToDeploy.Type {
	case "database":
		// For database services, no domain/subdomain is needed
		return deployDatabaseService(deployService, serviceToDeploy)
	case "llm":
		return deployDatabaseService(deployService, serviceToDeploy)
	case "github_repo":
		// For GitHub repos, we use the domain/subdomain pattern
		return deployGithubService(deployService, serviceToDeploy, baseDomain)
	default:
		// For backward compatibility or other service types, default to GitHub repo deployment
		fmt.Printf("Unspecified service type '%s', defaulting to GitHub repo deployment\n", serviceToDeploy.Type)
		return deployGithubService(deployService, serviceToDeploy, baseDomain)
	}
}

func DeployApplication(deployService *service.DeployService) error {
	services := deployService.DatabaseAdapter.GetServices()

	server := deployService.DatabaseAdapter.GetServer()

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

		err = deployService.DockerAdapter.RunRouterWithServer(server)
		if err != nil {
			return fmt.Errorf("error running Traefik router: %w", err)
		}
	} else {
		fmt.Println("Traefik router is already running")
	}

	if server.Domain == "" {
		return errors.New("server does not have domain")
	}

	for _, service := range services {
		fmt.Printf("Deploying service: %s\n", service.Name)
		err = deployOneService(deployService, service, server.Domain)
		if err != nil {
			service.Status = "Error"
			deployService.DatabaseAdapter.SaveService(service)
			return fmt.Errorf("error deploying service: %w", err)
		}
	}

	return nil
}
