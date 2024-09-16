package application

import (
	"fmt"
	"slices"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
)

func DeleteService(deployService *service.DeployService, serviceId string) error {
	service, err := deployService.DatabaseAdapter.GetServiceById(serviceId)

	if err != nil {
		fmt.Println(err)
	}

	projects := deployService.DatabaseAdapter.GetProjects()

	for _, p := range projects {
		for _, s := range p.Services {
			if s.Id == serviceId {
				service = &s
			}
		}
	}

	// TODO: move to a other method

	if service.DeployId != nil {
		return deleteServiceWithDeploy(deployService, service)
	} else {
		return deleteServiceWithoutDeploy(deployService, service)
	}
}

func deleteServiceWithDeploy(deployService *service.DeployService, s *domain.Service) error {
	deploy, err := deployService.DatabaseAdapter.GetDeployById(*s.DeployId)
	if err != nil {
		return err
	}

	server, err := deployService.DatabaseAdapter.GetServerById(deploy.ServerId)
	if err != nil {
		return err
	}

	deployService.DockerAdapter.ConnectClient(server)
	deployService.DockerAdapter.Delete(s.Name)

	err = deployService.DatabaseAdapter.DeleteServiceById(s.Id)
	if err != nil {
		return err
	}

	dEnvs := []dto.Env{}
	for _, dEnv := range deploy.Envs {
		if !slices.Contains(s.Envs, dEnv) {
			dEnvs = append(dEnvs, dEnv)
		}
	}

	deploy.Envs = dEnvs
	err = EditDeploy(deployService, dto.EditDeployDto{Id: deploy.Id, DeployOnCommit: deploy.DeployOnCommit, Envs: deploy.Envs, SubDomain: deploy.SubDomain})
	if err != nil {
		return err
	}

	return ReDeployApplication(deployService, deploy.Id)
}

func deleteServiceWithoutDeploy(deployService *service.DeployService, s *domain.Service) error {
	server := deployService.DockerAdapter.GetLocalHostServer()
	deployService.DockerAdapter.ConnectClient(server)
	deployService.DockerAdapter.Delete(s.Name)

	project, err := deployService.DatabaseAdapter.GetProjectById(*s.ProjectId)

	if err != nil {
		return err
	}

	deployService.FilesystemAdapter.RemoveEnvsFromDotEnvFile(project.Path, s.Envs)

	if err != nil {
		return err
	}

	return deployService.DatabaseAdapter.DeleteServiceById(s.Id)
}
