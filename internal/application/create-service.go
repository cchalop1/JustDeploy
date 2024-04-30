package application

import (
	"errors"
	"fmt"
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

func CreateService(deployService *service.DeployService, deployId string, createServiceDto dto.CreateServiceDto) error {
	deploy, err := deployService.DatabaseAdapter.GetDeployById(deployId)
	if err != nil {
		return err
	}
	server, err := deployService.DatabaseAdapter.GetServerById(deploy.ServerId)
	if err != nil {
		return err
	}

	fmt.Println(createServiceDto.ServiceName)

	service := database.ServicesConfig{}

	if createServiceDto.FromDockerCompose {
		// services, err := deployService.FilesystemAdapter.GetComposeConfigOfDeploy(deploy.PathToSource)

		// if err != nil {
		// 	return err
		// }

		// for _, s := range services {
		// 	if s.Name == createServiceDto.ServiceName {
		// 		// service = s
		// 	}
		// }
		return errors.New("Not implemented")

	} else {
		service, err = database.GetServiceByName(createServiceDto.ServiceName)
	}
	fmt.Println(service)

	if err != nil {
		return err
	}

	deployService.DockerAdapter.ConnectClient(server)

	deployService.DockerAdapter.PullImage(service.Config.Image)

	envs := generateEnvsForService(service.Env)

	// replace in services config all the envs with the new values
	replaceEnvForServicesConfig(&service, envs)

	containerHostname := strings.ToLower(service.Name) + "-db-" + deployId + "-" + utils.GenerateRandomPassword(5)

	deployService.DockerAdapter.RunService(service, containerHostname)

	envs = append(envs, dto.Env{Name: strings.ToUpper(service.Name) + "_HOSTNAME", Value: containerHostname})

	// TODO: add volume
	domainService := domain.Service{
		Id:          utils.GenerateRandomPassword(5),
		DeployId:    deployId,
		Name:        containerHostname,
		Envs:        envs,
		VolumsNames: []string{},
		Status:      "Runing",
		ImageName:   service.Config.Image,
	}

	deployService.DatabaseAdapter.SaveService(domainService)

	EditDeploy(deployService, dto.EditDeployDto{Id: deploy.Id, Envs: append(envs, deploy.Envs...), SubDomain: deploy.SubDomain, DeployOnCommit: deploy.DeployOnCommit})

	ReDeployApplication(deployService, deploy.Id)

	return nil
}
