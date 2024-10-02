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

	project, err := deployService.DatabaseAdapter.GetProjectById(deployProjectDto.ProjectId)
	if err != nil {
		fmt.Println("Error fetching project by ID:", err)
		return domain.Deploy{}, err
	}

	server, err := deployService.DatabaseAdapter.GetServerById(deployProjectDto.ServerId)
	if err != nil {
		return domain.Deploy{}, err
	}

	err = deployService.DockerAdapter.ConnectClient(server)
	if err != nil {
		return domain.Deploy{}, err
	}

	exposeServiceId := ""
	for _, service := range project.Services {
		if service.IsDevContainer {
			exposeServiceId = service.Id
			fmt.Println("Service to expose found:", service.Name)
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
		return domain.Deploy{}, err
	}

	// Build all services
	for _, service := range project.Services {
		err = pullAndBuildService(deployService, service, server)
		if err != nil {
			return domain.Deploy{}, err
		}
	}

	containersConfig := []container.Config{}

	// Configure all services
	for _, service := range project.Services {
		config := deployService.DockerAdapter.ConfigContainer(service)
		if service.Id == exposeServiceId {
			deployService.DockerAdapter.ExposeContainer(&config, adapter.ExposeContainerParams{
				IsTls:  true,
				Domain: server.Domain,
				Port:   service.ExposePort,
			})
		}
		containersConfig = append(containersConfig, config)
	}

	// Run all services
	networkName := project.Name + "_network"
	for _, config := range containersConfig {
		err = deployService.DockerAdapter.RunImage(config, networkName)
		if err != nil {
			return domain.Deploy{}, err
		}
	}

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
