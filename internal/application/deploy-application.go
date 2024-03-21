package application

import (
	"fmt"
	"path/filepath"

	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/internal/utils"
)

func runApplication(deployService *service.DeployService, deploy domain.Deploy, domain string) {
	deployService.DockerAdapter.BuildImage(deploy.Name, deploy.PathToSource)
	deployService.DockerAdapter.PullTreafikImage()
	deployService.DockerAdapter.RunRouter()
	deployService.DockerAdapter.RunImage(deploy, domain)

	if deploy.EnableTls {
		deploy.Url = "https://" + domain
	} else {
		deploy.Url = "http://" + domain
	}

	deploy.Status = "Runing"
	deployService.DatabaseAdapter.UpdateDeploy(deploy)
}

func DeployApplication(deployService *service.DeployService, newDeploy dto.NewDeployDto) error {
	server, err := deployService.DatabaseAdapter.GetServerById(newDeploy.ServerId)
	if err != nil {
		return err
	}

	pathToDir, err := filepath.Abs(newDeploy.PathToSource)

	if err != nil {
		return err
	}

	pathToDir = adapter.NewFilesystemAdapter().CleanPath(pathToDir)

	err = deployService.DockerAdapter.ConnectClient(server)

	if err != nil {
		return err
	}

	deploy := domain.Deploy{
		Id:           utils.GenerateRandomPassword(5),
		Name:         newDeploy.Name,
		ServerId:     newDeploy.ServerId,
		PathToSource: pathToDir,
		Status:       "Installing",
	}

	err = deployService.DatabaseAdapter.SaveDeploy(deploy)
	fmt.Println(err)

	runApplication(deployService, deploy, server.Domain)

	return nil
}

func ReDeployApplicationRun(deployService *service.DeployService, deploy domain.Deploy) error {
	server, err := deployService.DatabaseAdapter.GetServerById(deploy.ServerId)
	if err != nil {
		return err
	}

	err = deployService.DockerAdapter.ConnectClient(server)

	if err != nil {
		return err
	}

	deploy.Status = "Installing"

	err = deployService.DatabaseAdapter.SaveDeploy(deploy)
	fmt.Println(err)
	runApplication(deployService, deploy, server.Domain)
	return nil

}
