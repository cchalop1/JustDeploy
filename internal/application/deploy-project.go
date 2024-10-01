package application

import (
	"fmt"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
)

func DeployProject(deployService *service.DeployService, deployProjectDto dto.DeployProjectDto) error {
	project, err := deployService.DatabaseAdapter.GetProjectById(deployProjectDto.ProjectId)

	if err != nil {
		return err
	}

	server, err := deployService.DatabaseAdapter.GetServerById(deployProjectDto.ServerId)

	if err != nil {
		return err
	}

	err = deployService.DockerAdapter.ConnectClient(server)

	if err != nil {
		return err
	}

	fmt.Println("Deploy project ", project)

	return nil
}
