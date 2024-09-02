package application

import (
	"strings"

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

func createServiceLinkToDeploy(deployService *service.DeployService, createServiceDto dto.CreateServiceDto) error {
	deploy, err := deployService.DatabaseAdapter.GetDeployById(*createServiceDto.DeployId)
	if err != nil {
		return err
	}

	server, err := deployService.DatabaseAdapter.GetServerById(deploy.ServerId)
	if err != nil {
		return err
	}

	// extract service from deploy
	service := database.ServicesConfig{}

	if createServiceDto.FromDockerCompose {
		services, err := deployService.FilesystemAdapter.GetComposeConfigOfDeploy(deploy.PathToSource)

		if err != nil {
			return err
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
		return err
	}

	deployService.DockerAdapter.ConnectClient(server)

	deployService.DockerAdapter.PullImage(service.Config.Image)

	envs := generateEnvsForService(service.Env)

	// replace in services config all the envs with the new values
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

	return nil
}

func createServiceOnly(deployService *service.DeployService, serviceName string) error {
	server := deployService.DockerAdapter.GetLocalHostServer()

	service, err := database.GetServiceByName(serviceName)

	if err != nil {
		return err
	}

	deployService.DockerAdapter.ConnectClient(server)

	deployService.DockerAdapter.PullImage(service.Config.Image)

	envs := generateEnvsForService(service.Env)

	// replace in services config all the envs with the new values
	replaceEnvForServicesConfig(&service, envs)

	containerHostname := generateContainerHostname(service.Name, nil)

	deployService.DockerAdapter.RunService(service, containerHostname)

	envs = append(envs, dto.Env{Name: "localhost", Value: containerHostname})

	domainService := domain.Service{
		Id:          utils.GenerateRandomPassword(5),
		DeployId:    nil,
		Name:        containerHostname,
		Envs:        envs,
		VolumsNames: []string{},
		Status:      "Runing",
		ImageName:   service.Config.Image,
		ImageUrl:    service.Icon,
	}

	deployService.DatabaseAdapter.SaveService(domainService)

	return nil
}

func CreateService(deployService *service.DeployService, createServiceDto dto.CreateServiceDto) error {
	if createServiceDto.DeployId == nil {
		return createServiceOnly(deployService, createServiceDto.ServiceName)
	} else {
		return createServiceLinkToDeploy(deployService, createServiceDto)
	}
}
