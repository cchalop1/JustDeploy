package application

import (
	"fmt"
	"strings"

	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/internal/utils"
	"github.com/docker/docker/api/types/container"
)

func DeployProject(deployService *service.DeployService, deployProjectDto dto.DeployProjectDto) (domain.Deploy, error) {
	fmt.Println("Starting deployment process for project ID:", deployProjectDto.ProjectId)

	project, err := deployService.DatabaseAdapter.GetProjectById(deployProjectDto.ProjectId)
	if err != nil {
		fmt.Println("Error fetching project by ID:", err)
		return domain.Deploy{}, err
	}

	fmt.Println("Fetched project:", project.Name)

	server, err := deployService.DatabaseAdapter.GetServerById(deployProjectDto.ServerId)
	if err != nil {
		fmt.Println("Error fetching server by ID:", err)
		return domain.Deploy{}, err
	}
	fmt.Println("Fetched server:", server.Name)

	err = deployService.DockerAdapter.ConnectClient(server)
	if err != nil {
		fmt.Println("Error connecting Docker client:", err)
		return domain.Deploy{}, err
	}
	fmt.Println("Connected to Docker client")

	exposeServiceId := ""
	for _, service := range project.Services {
		if service.IsDevContainer {
			exposeServiceId = service.Id
			fmt.Println("Service to expose found:", service.HostName)
			break
		}
	}

	deploy := domain.Deploy{
		Id:              utils.GenerateRandomPassword(5),
		ServerId:        server.Id,
		ProjectId:       project.Id,
		EnableTls:       deployProjectDto.IsTLSDomain,
		Email:           server.Email,
		ExposeServiceId: exposeServiceId,
		ServicesDeploy:  project.Services,
	}

	err = deployService.DatabaseAdapter.SaveDeploy(deploy)
	if err != nil {
		fmt.Println("Error saving deploy:", err)
		return domain.Deploy{}, err
	}
	fmt.Println("Deploy saved with ID:", deploy.Id)

	// Build all services
	for _, service := range project.Services {
		fmt.Println("Building service:", service.HostName)
		err = pullAndBuildService(deployService, service, server)
		if err != nil {
			fmt.Println("Error building service:", service.HostName, err)
			return domain.Deploy{}, err
		}
		fmt.Println("Service built:", service.HostName)
	}

	containersConfig := []container.Config{}

	// TODO: regeneate envs and hostname
	for i, service := range project.Services {
		envsFiltered := []dto.Env{}
		for i, env := range service.Envs {
			if !strings.Contains(env.Name, "HOSTNAME") {
				envsFiltered = append(envsFiltered, service.Envs[i])
			}
		}
		project.Services[i].Envs = envsFiltered
	}

	hostNameEnvs := []dto.Env{}

	for _, service := range project.Services {
		if !service.IsDevContainer {
			hostNameEnvs = append(hostNameEnvs, dto.Env{
				//TODO: get the real name of the service
				Name:  strings.ToUpper(service.Name) + "_HOSTNAME",
				Value: service.GetDockerName(),
			})
		}
	}

	for i, service := range project.Services {
		if service.IsDevContainer {
			project.Services[i].Envs = append(service.Envs, hostNameEnvs...)
		}
	}

	// Configure all services
	for _, service := range project.Services {
		fmt.Println("Configuring service:", service.HostName)
		config := deployService.DockerAdapter.ConfigContainer(service)
		if service.Id == exposeServiceId {
			fmt.Println("Exposing service:", service.HostName)
			deployService.DockerAdapter.ExposeContainer(&config, adapter.ExposeContainerParams{
				IsTls:  true,
				Domain: server.Domain,
				Port:   service.ExposePort,
			})
		}
		containersConfig = append(containersConfig, config)
		fmt.Println("Service configured:", service.HostName)
	}

	// Run all services
	networkName := project.Name + "_network"
	for _, config := range containersConfig {
		fmt.Println("Running service with config:", config.Image)
		err = deployService.DockerAdapter.RunImage(config, networkName)
		if err != nil {
			fmt.Println("Error running service with config:", config.Image, err)
			return domain.Deploy{}, err
		}
		fmt.Println("Service running with config:", config.Image)
	}

	project.ServerId = &server.Id

	err = deployService.DatabaseAdapter.SaveProject(*project)
	if err != nil {
		fmt.Println("Error saving project:", err)
		return domain.Deploy{}, err
	}

	fmt.Println("Deployment process completed for project ID:", deployProjectDto.ProjectId)
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
