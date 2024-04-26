package application

import (
	"errors"
	"strings"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/internal/utils"
)

func findServiceByName(deployService *service.DeployService, serviceName string) (dto.ServiceDto, error) {
	services := deployService.FilesystemAdapter.GetServicesListConfig()

	for _, service := range services {
		if service.Name == serviceName {
			return service, nil
		}
	}

	return dto.ServiceDto{}, errors.New("service not found")
}

func generateEnvs(service dto.ServiceDto) []dto.Env {
	envs := []dto.Env{}

	for _, env := range service.Envs {
		envs = append(envs, dto.Env{Name: env, Secret: strings.ToLower(env)})
	}

	for _, secret := range service.Secrets {
		envs = append(envs, dto.Env{Name: secret, Secret: utils.GenerateRandomPassword(12)})
	}

	return envs
}

func replaceConnectUrl(connectUrl string, envs []dto.Env) string {
	for _, env := range envs {
		connectUrl = strings.Replace(connectUrl, env.Name, env.Secret, -1)
	}
	return connectUrl
}

func CreateService(deployService *service.DeployService, deployId string, createServiceDto dto.CreateServiceDto) error {
	deploy, err := deployService.DatabaseAdapter.GetDeployById(deployId)
	if err != nil {
		return err
	}
	server, err := deployService.DatabaseAdapter.GetServerById(deploy.ServerId)
	if err != nil {
		return err
	}

	service := dto.ServiceDto{}

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
		service, err = findServiceByName(deployService, createServiceDto.ServiceName)
	}

	if err != nil {
		return err
	}

	deployService.DockerAdapter.ConnectClient(server)
	deployService.DockerAdapter.PullImage(service.Image)

	envs := generateEnvs(service)

	containerHostname := strings.ToLower(service.Name) + "-db-" + deployId + "-" + utils.GenerateRandomPassword(5)

	deployService.DockerAdapter.RunService(service, envs, containerHostname)

	envs = append(envs, dto.Env{Name: strings.ToUpper(service.Name) + "_HOSTNAME", Secret: containerHostname})

	// TODO: add volume

	domainService := domain.Service{
		Id:          utils.GenerateRandomPassword(5),
		DeployId:    deployId,
		Name:        containerHostname,
		Envs:        envs,
		VolumsNames: []string{},
		Status:      "Runing",
		ImageName:   service.Image,
	}

	deployService.DatabaseAdapter.SaveService(domainService)

	EditDeploy(deployService, dto.EditDeployDto{Id: deploy.Id, Envs: append(envs, deploy.Envs...), SubDomain: deploy.SubDomain, DeployOnCommit: deploy.DeployOnCommit})

	ReDeployApplication(deployService, deploy.Id)

	return nil
}
