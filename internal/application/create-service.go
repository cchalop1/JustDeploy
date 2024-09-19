package application

import (
	"errors"
	"fmt"
	"strings"

	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/adapter/database"
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/internal/utils"
)

func generateEnvsForService(serviceEnvs []dto.Env) []dto.Env {
	envs := []dto.Env{}

	for idx := range serviceEnvs {
		if serviceEnvs[idx].IsSecret {
			envs = append(envs, dto.Env{Name: serviceEnvs[idx].Name, Value: utils.GenerateRandomPassword(12), IsSecret: true})
		} else {
			envs = append(envs, dto.Env{Name: serviceEnvs[idx].Name, Value: strings.ToLower(serviceEnvs[idx].Name), IsSecret: true})
		}
	}

	return envs
}

func replaceEnvForServicesConfig(service *database.ServicesConfig, serviceEnvs []dto.Env) {
	service.Env = serviceEnvs
	for i, envStr := range service.Config.Env {
		for _, env := range serviceEnvs {
			envStr = strings.ReplaceAll(envStr, "$"+env.Name, env.Value)
		}
		service.Config.Env[i] = envStr
	}

	for _, env := range serviceEnvs {
		for idx, cmd := range service.Config.Cmd {
			cmd := strings.ReplaceAll(cmd, "$"+env.Name, env.Value)
			service.Config.Cmd[idx] = cmd
		}
	}
}

func generateContainerHostname(serviceName string, deployId *string) string {
	serviceName += "-"

	if deployId != nil {
		serviceName += *deployId + "-"
	}

	return serviceName + utils.GenerateRandomPassword(5)
}

func getPortsForService(service database.ServicesConfig) string {
	return adapter.FindOpenLocalPort(service.DefaultPort)
}

func createServiceLinkToDeploy(deployService *service.DeployService, createServiceDto dto.CreateServiceDto) (domain.Service, error) {
	deploy, err := deployService.DatabaseAdapter.GetDeployById(*createServiceDto.DeployId)
	if err != nil {
		return domain.Service{}, err
	}

	server, err := deployService.DatabaseAdapter.GetServerById(deploy.ServerId)
	if err != nil {
		return domain.Service{}, err
	}

	// extract service from deploy
	service := database.ServicesConfig{}

	if createServiceDto.FromDockerCompose {
		services, err := deployService.FilesystemAdapter.GetComposeConfigOfDeploy(deploy.PathToSource)

		if err != nil {
			return domain.Service{}, err
		}

		for _, s := range services {
			if s.Name == createServiceDto.ServiceName {
				service = s
			}
		}

	} else {
		service, err = database.GetServiceByName(createServiceDto.ServiceName)
	}

	if err != nil {
		return domain.Service{}, err
	}

	deployService.DockerAdapter.ConnectClient(server)

	deployService.DockerAdapter.PullImage(service.Config.Image)

	envs := generateEnvsForService(service.Env)

	replaceEnvForServicesConfig(&service, envs)

	containerHostname := generateContainerHostname(service.Name, createServiceDto.DeployId)

	deployService.DockerAdapter.RunServiceWithDeploy(service, containerHostname)

	envs = append(envs, dto.Env{Name: strings.ToUpper(service.Name) + "_HOSTNAME", Value: containerHostname})

	// TODO: add volume
	domainService := domain.Service{
		Id:          utils.GenerateRandomPassword(5),
		DeployId:    createServiceDto.DeployId,
		Name:        containerHostname,
		Envs:        envs,
		VolumsNames: []string{},
		Status:      "Runing",
		ImageName:   service.Config.Image,
		ImageUrl:    service.Icon,
	}

	deployService.DatabaseAdapter.SaveService(domainService)

	EditDeploy(deployService, dto.EditDeployDto{Id: deploy.Id, Envs: append(envs, deploy.Envs...), SubDomain: deploy.SubDomain, DeployOnCommit: deploy.DeployOnCommit})
	ReDeployApplication(deployService, deploy.Id)

	return domainService, nil
}

func createDevContainerService(deployService *service.DeployService, createServiceDto dto.CreateServiceDto) (domain.Service, error) {
	domainService := domain.Service{
		Id:   utils.GenerateRandomPassword(5),
		Name: deployService.FilesystemAdapter.GetFolderName(*createServiceDto.Path),
		Envs: []dto.Env{
			// TODO: find a avalaible port
			{Name: "PORT", Value: "9999"},
		},
		VolumsNames:    []string{},
		Status:         "Runing",
		Host:           "localhost",
		ProjectId:      createServiceDto.ProjectId,
		IsDevContainer: true,
		CurrentPath:    createServiceDto.Path,
	}
	err := deployService.DatabaseAdapter.SaveServiceByProjectId(domainService)
	return domainService, err
}

func createServiceForProject(deployService *service.DeployService, createServiceDto dto.CreateServiceDto) (domain.Service, error) {
	if createServiceDto.ProjectId == nil {
		return domain.Service{}, errors.New("ProjectId is required")
	}

	fmt.Println("createServiceDto.Path", createServiceDto.Path)
	if createServiceDto.Path != nil {
		return createDevContainerService(deployService, createServiceDto)
	}

	service, err := database.GetServiceByName(createServiceDto.ServiceName)

	if err != nil {
		return domain.Service{}, err
	}

	server := deployService.DockerAdapter.GetLocalHostServer()
	deployService.DockerAdapter.ConnectClient(server)

	deployService.DockerAdapter.PullImage(service.Config.Image)

	envs := generateEnvsForService(service.Env)

	replaceEnvForServicesConfig(&service, envs)

	containerHostname := generateContainerHostname(service.Name, nil)

	exposedPort := getPortsForService(service)

	deployService.DockerAdapter.RunService(service, exposedPort, containerHostname)

	envs = append(envs, dto.Env{Name: strings.ToUpper(service.Name) + "_HOSTNAME", Value: "localhost"})

	envs = append(envs, dto.Env{Name: strings.ToUpper(service.Name) + "_PORT", Value: exposedPort})

	domainService := domain.Service{
		Id:          utils.GenerateRandomPassword(5),
		DeployId:    nil,
		ProjectId:   createServiceDto.ProjectId,
		Name:        containerHostname,
		Envs:        envs,
		VolumsNames: []string{},
		Status:      "Runing",
		ImageName:   service.Config.Image,
		ImageUrl:    service.Icon,
		ExposePort:  &exposedPort,
	}

	project, err := deployService.DatabaseAdapter.GetProjectById(*createServiceDto.ProjectId)

	if err != nil {
		return domain.Service{}, err
	}

	project.Services = append(project.Services, domainService)

	// TODO: add envs to devContainers service in project
	for i, service := range project.Services {
		if project.Services[i].IsDevContainer {
			project.Services[i].Envs = append(service.Envs, domainService.Envs...)
		}
	}

	deployService.DatabaseAdapter.SaveProject(*project)

	deployService.FilesystemAdapter.GenerateDotEnvFile(project.Path, envs)

	return domainService, nil
}

func CreateService(deployService *service.DeployService, createServiceDto dto.CreateServiceDto) (domain.Service, error) {
	if createServiceDto.DeployId == nil {
		return createServiceForProject(deployService, createServiceDto)
	} else {
		return createServiceLinkToDeploy(deployService, createServiceDto)
	}
}
