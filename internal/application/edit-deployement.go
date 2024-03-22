package application

import (
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
)

func EditDeploy(deployService *service.DeployService, editDeployDto dto.EditDeployDto) error {
	deploy, err := deployService.DatabaseAdapter.GetDeployById(editDeployDto.Id)
	if err != nil {
		return err
	}
	deploy.Envs = editDeployDto.Envs
	deploy.SubDomain = editDeployDto.SubDomain
	deploy.DeployOnCommit = editDeployDto.DeployOnCommit

	if deploy.DeployOnCommit {
		deployService.FilesystemAdapter.CreateGitPostCommitHooks(deploy)
	} else {
		deployService.FilesystemAdapter.DeleteGitPostCommitHooks(deploy)
	}
	deployService.DatabaseAdapter.UpdateDeploy(deploy)
	return nil
}
