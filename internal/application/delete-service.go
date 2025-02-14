package application

import (
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
)

func DeleteService(deployService *service.DeployService, serviceId string) error {
	s, err := deployService.DatabaseAdapter.GetServiceById(serviceId)

	if err != nil {
		return err
	}

	if s.Status != "stoped" || s.Status != "ready_to_deploy" {
		// TODO: return error message
	}

	return deployService.DatabaseAdapter.DeleteServiceById(s.Id)

	// return deleteServiceWithoutDeploy(deployService, p, service)
}

func deleteServiceWithoutDeploy(deployService *service.DeployService, project *domain.Project, s *domain.Service) error {
	server := deployService.DockerAdapter.GetLocalHostServer()
	deployService.DockerAdapter.ConnectClient(server)
	deployService.DockerAdapter.Delete(s.HostName)

	removeEnvsFromProject(project, s.Envs)

	err := deployService.DatabaseAdapter.SaveProject(*project)

	if err != nil {
		return err
	}

	return deployService.DatabaseAdapter.DeleteServiceById(s.Id)
}

func removeEnvsFromProject(project *domain.Project, envsToRemove []dto.Env) {
	for i, _ := range project.Services {
		if project.Services[i].IsDevContainer {
			filteredEnvs := []dto.Env{}
			for _, env := range project.Services[i].Envs {
				if !envIsContained(env, envsToRemove) {
					filteredEnvs = append(filteredEnvs, env)
				}
			}
			project.Services[i].Envs = filteredEnvs
		}
	}
}

func envIsContained(env dto.Env, envs []dto.Env) bool {
	for _, e := range envs {
		if e.Equals(env) {
			return true
		}
	}
	return false
}
