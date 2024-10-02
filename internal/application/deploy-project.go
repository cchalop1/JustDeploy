package application

import (
	"fmt"

	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/internal/utils"
	"github.com/docker/docker/api/types/container"
)

func DeployProject(deployService *service.DeployService, deployProjectDto dto.DeployProjectDto) (domain.Deploy, error) {
	fmt.Println("Starting DeployProject function")

	fmt.Println("Fetching project by ID:", deployProjectDto.ProjectId)
	project, err := deployService.DatabaseAdapter.GetProjectById(deployProjectDto.ProjectId)
	if err != nil {
		fmt.Println("Error fetching project by ID:", err)
		return domain.Deploy{}, err
	}
	fmt.Println("Project found:", project.Name)

	fmt.Println("Fetching server by ID:", deployProjectDto.ServerId)
	server, err := deployService.DatabaseAdapter.GetServerById(deployProjectDto.ServerId)
	if err != nil {
		fmt.Println("Error fetching server by ID:", err)
		return domain.Deploy{}, err
	}
	fmt.Println("Server found:", server.Domain)

	fmt.Println("Connecting Docker client to server:", server.Domain)
	err = deployService.DockerAdapter.ConnectClient(server)
	if err != nil {
		fmt.Println("Error connecting to Docker client:", err)
		return domain.Deploy{}, err
	}

	fmt.Println("Determining service to expose")
	exposeServiceId := ""
	for _, service := range project.Services {
		if service.IsDevContainer {
			exposeServiceId = service.Id
			fmt.Println("Service to expose found:", service.Name)
			break
		}
	}

	fmt.Println("Creating deploy object")
	deploy := domain.Deploy{
		Id:              utils.GenerateRandomPassword(5),
		ServerId:        server.Id,
		ProjectId:       project.Id,
		EnableTls:       deployProjectDto.IsTLSDomain,
		Email:           server.Email,
		ExposeServiceId: exposeServiceId,
		ServicesDeploy:  project.Services,
	}

	fmt.Println("Saving deploy object to database")
	err = deployService.DatabaseAdapter.SaveDeploy(deploy)
	if err != nil {
		fmt.Println("Error saving deploy object:", err)
		return domain.Deploy{}, err
	}

	// Build all services
	fmt.Println("Building all services")
	for _, service := range project.Services {
		fmt.Println("Building service:", service.Name)
		err = pullAndBuildService(deployService, service, server)
		if err != nil {
			fmt.Println("Error building service:", service.Name, err)
			return domain.Deploy{}, err
		}
	}

	containersConfig := []container.Config{}

	// Configure all services
	fmt.Println("Configuring all services")
	for _, service := range project.Services {
		fmt.Println("Configuring service:", service.Name)
		config := deployService.DockerAdapter.ConfigContainer(service)
		if service.Id == exposeServiceId {
			fmt.Println("Exposing container for service:", service.Name)
			deployService.DockerAdapter.ExposeContainer(&config, adapter.ExposeContainerParams{
				IsTls:  true,
				Domain: server.Domain,
			})
		}
		containersConfig = append(containersConfig, config)
	}

	// Run all services
	networkName := project.Name + "_network"
	fmt.Println("Running all services on network:", networkName)
	for _, config := range containersConfig {
		fmt.Println("Starting service with image:", config.Image)
		err = deployService.DockerAdapter.RunImage(config, networkName)
		if err != nil {
			fmt.Println("Error starting service with image:", config.Image, err)
			return domain.Deploy{}, err
		}
	}

	fmt.Println("DeployProject function completed successfully")
	return deploy, nil
}

func pullAndBuildService(deployService *service.DeployService, service domain.Service, server domain.Server) error {
	var err error

	if service.IsDevContainer {
		isDockerFile := deployService.FilesystemAdapter.FindDockerFile(service.CurrentPath)
		if isDockerFile {
			err = deployService.DockerAdapter.BuildImage(service)
		} else {
			err = deployService.DockerAdapter.BuildNixpacksImage(server, service)
		}
	} else {
		err = deployService.DockerAdapter.PullImage(service.ImageName)
	}

	return err
}
