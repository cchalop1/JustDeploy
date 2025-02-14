package application

// import (
// 	"errors"
// 	"strings"

// 	"cchalop1.com/deploy/internal/adapter"
// 	"cchalop1.com/deploy/internal/adapter/database"
// 	"cchalop1.com/deploy/internal/api/dto"
// 	"cchalop1.com/deploy/internal/api/service"
// 	"cchalop1.com/deploy/internal/domain"
// 	"cchalop1.com/deploy/internal/utils"
// )

// func generateEnvsForService(serviceEnvs []dto.Env) []dto.Env {
// 	envs := []dto.Env{}

// 	for idx := range serviceEnvs {
// 		if serviceEnvs[idx].IsSecret {
// 			envs = append(envs, dto.Env{Name: serviceEnvs[idx].Name, Value: utils.GenerateRandomPassword(12), IsSecret: true})
// 		} else {
// 			envs = append(envs, dto.Env{Name: serviceEnvs[idx].Name, Value: strings.ToLower(serviceEnvs[idx].Name), IsSecret: true})
// 		}
// 	}

// 	return envs
// }

// func replaceEnvForServicesConfig(service *database.ServicesConfig, serviceEnvs []dto.Env) {
// 	service.Env = serviceEnvs
// 	for i, envStr := range service.Config.Env {
// 		for _, env := range serviceEnvs {
// 			envStr = strings.ReplaceAll(envStr, "$"+env.Name, env.Value)
// 		}
// 		service.Config.Env[i] = envStr
// 	}

// 	for _, env := range serviceEnvs {
// 		for idx, cmd := range service.Config.Cmd {
// 			cmd := strings.ReplaceAll(cmd, "$"+env.Name, env.Value)
// 			service.Config.Cmd[idx] = cmd
// 		}
// 	}
// }

// func generateContainerHostname(serviceName string, deployId *string) string {
// 	serviceName += "-"

// 	if deployId != nil {
// 		serviceName += *deployId + "-"
// 	}

// 	return serviceName + utils.GenerateRandomPassword(5)
// }

// func getPortsForService(service database.ServicesConfig) string {
// 	return adapter.FindOpenLocalPort(service.DefaultPort)
// }

// func createDevContainerService(deployService *service.DeployService, createServiceDto dto.CreateServiceDto) (domain.Service, error) {
// 	exposedPort := "9999"
// 	envs := deployService.FilesystemAdapter.LoadEnvsFromFileSystem(*createServiceDto.Path)

// 	for idx := range envs {
// 		if envs[idx].Name == "PORT" {
// 			exposedPort = envs[idx].Value
// 		}
// 	}
// 	// remove PORT env form envs
// 	for idx, env := range envs {
// 		if env.Name == "PORT" {
// 			envs = append(envs[:idx], envs[idx+1:]...)
// 		}
// 	}

// 	envs = append(envs, dto.Env{Name: "PORT", Value: exposedPort})

// 	domainService := domain.Service{
// 		Id:             utils.GenerateRandomPassword(5),
// 		HostName:       deployService.FilesystemAdapter.GetFolderName(*createServiceDto.Path),
// 		Envs:           envs,
// 		VolumsNames:    []string{},
// 		Status:         "Runing",
// 		Host:           "localhost",
// 		IsDevContainer: true,
// 		CurrentPath:    *createServiceDto.Path,
// 		ExposePort:     exposedPort,
// 	}
// 	err := deployService.DatabaseAdapter.SaveServiceByProjectId(*createServiceDto.ProjectId, domainService)
// 	return domainService, err
// }

// func createServiceForProject(deployService *service.DeployService, createServiceDto dto.CreateServiceDto) (domain.Service, error) {
// 	if createServiceDto.ProjectId == nil {
// 		return domain.Service{}, errors.New("ProjectId is required")
// 	}

// 	if createServiceDto.Path != nil {
// 		return createDevContainerService(deployService, createServiceDto)
// 	}

// 	service, err := database.GetServiceByName(createServiceDto.ServiceName)

// 	if err != nil {
// 		return domain.Service{}, err
// 	}

// 	server := deployService.DockerAdapter.GetLocalHostServer()
// 	deployService.DockerAdapter.ConnectClient(server)

// 	deployService.DockerAdapter.PullImage(service.Config.Image)

// 	envs := generateEnvsForService(service.Env)

// 	replaceEnvForServicesConfig(&service, envs)

// 	containerHostname := generateContainerHostname(service.Name, nil)

// 	exposedPort := getPortsForService(service)

// 	deployService.DockerAdapter.RunService(service, exposedPort, containerHostname)

// 	envs = append(envs, dto.Env{Name: strings.ToUpper(service.Name) + "_HOSTNAME", Value: "localhost"})

// 	envs = append(envs, dto.Env{Name: strings.ToUpper(service.Name) + "_PORT", Value: exposedPort})

// 	domainService := domain.Service{
// 		Id:          utils.GenerateRandomPassword(5),
// 		HostName:    containerHostname,
// 		Envs:        envs,
// 		VolumsNames: []string{},
// 		Status:      "Runing",
// 		Name:        strings.ToLower(service.Name),
// 		ImageName:   service.Config.Image,
// 		ImageUrl:    service.Icon,
// 		ExposePort:  exposedPort,
// 	}

// 	project, err := deployService.DatabaseAdapter.GetProjectById(*createServiceDto.ProjectId)

// 	if err != nil {
// 		return domain.Service{}, err
// 	}

// 	project.Services = append(project.Services, domainService)

// 	addEnvsToProject(project, envs)

// 	deployService.DatabaseAdapter.SaveProject(*project)

// 	deployService.FilesystemAdapter.GenerateDotEnvFile(project)

// 	return domainService, nil
// }

// func CreateService(deployService *service.DeployService, createServiceDto dto.CreateServiceDto) (domain.Service, error) {
// 	return createServiceForProject(deployService, createServiceDto)
// }

// func addEnvsToProject(project *domain.Project, envs []dto.Env) {
// 	for i, service := range project.Services {
// 		if project.Services[i].IsDevContainer {
// 			project.Services[i].Envs = append(service.Envs, envs...)
// 		}
// 	}
// }
