package application

import (
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
)

func EditDeployement(deployService *service.DeployService, containerName string, editDeployDto dto.EditDeployementDto) {
	if editDeployDto.DeployOnCommit {
		deployService.FilesystemAdapter.CreateGitPostCommitHooks(deployService.DeployConfig.AppConfig.PathToSource)
	} else {
		deployService.FilesystemAdapter.DeleteGitPostCommitHooks(deployService.DeployConfig.AppConfig.PathToSource)
	}
	deployService.DeployConfig.AppConfig.DeployOnCommit = editDeployDto.DeployOnCommit
	deployService.DatabaseAdapter.SaveState(*deployService.DeployConfig)
}
