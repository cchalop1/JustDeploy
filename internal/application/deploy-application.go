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

func runApplication(deployService *service.DeployService, deploy *domain.Deploy, domain string) {
	appUrl := ""

	if deploy.SubDomain != "" {
		appUrl += deploy.SubDomain + "."
	}

	appUrl += domain

	deployService.DockerAdapter.BuildImage(deploy)
	deployService.DockerAdapter.PullTreafikImage()
	deployService.DockerAdapter.RunRouter()
	deployService.DockerAdapter.RunImage(deploy, appUrl)

	if deploy.EnableTls {
		deploy.Url = "https://" + appUrl
	} else {
		deploy.Url = "http://" + appUrl
	}

	fmt.Println(deploy.Url)
	deploy.Status = "Runing"
	deployService.DatabaseAdapter.UpdateDeploy(*deploy)
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

	isFolder := adapter.NewFilesystemAdapter().IsFolder(pathToDir)
	DockerFileName := "Dockerfile"
	fmt.Println(isFolder)
	if !isFolder {
		DockerFileName = adapter.NewFilesystemAdapter().BaseDir(pathToDir)
		pathToDir = adapter.NewFilesystemAdapter().GetDir(pathToDir)
	}

	deploy := domain.Deploy{
		Id:             utils.GenerateRandomPassword(5),
		Name:           newDeploy.Name,
		ServerId:       newDeploy.ServerId,
		PathToSource:   pathToDir,
		Status:         "Installing",
		EnableTls:      newDeploy.EnableTls,
		Email:          newDeploy.Email,
		Envs:           newDeploy.Envs,
		SubDomain:      newDeploy.Name,
		DockerFileName: DockerFileName,
	}

	err = deployService.DatabaseAdapter.SaveDeploy(deploy)
	fmt.Println(err)

	runApplication(deployService, &deploy, server.Domain)

	return nil
}

func ReDeployApplicationRun(deployService *service.DeployService, deploy *domain.Deploy) error {
	server, err := deployService.DatabaseAdapter.GetServerById(deploy.ServerId)
	if err != nil {
		return err
	}

	err = deployService.DockerAdapter.ConnectClient(server)

	if err != nil {
		return err
	}

	deploy.Status = "Installing"

	err = deployService.DatabaseAdapter.UpdateDeploy(*deploy)
	fmt.Println(err)
	runApplication(deployService, deploy, server.Domain)
	return nil

}
