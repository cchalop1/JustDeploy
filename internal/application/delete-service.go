package application

import (
	"slices"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
)

func DeleteService(deployService *service.DeployService, serviceId string) error {
	s, err := deployService.DatabaseAdapter.GetServiceById(serviceId)
	if err != nil {
		return err
	}

	deploy, err := deployService.DatabaseAdapter.GetDeployById(s.DeployId)
	if err != nil {
		return err
	}

	server, err := deployService.DatabaseAdapter.GetServerById(deploy.ServerId)
	if err != nil {
		return err
	}

	deployService.DockerAdapter.ConnectClient(server)
	deployService.DockerAdapter.Delete(s.Name, false)

	err = deployService.DatabaseAdapter.DeleteServiceById(serviceId)
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
	err = ReDeployApplication(deployService, deploy.Id)

	if err != nil {
		return err
	}

	return err
}
