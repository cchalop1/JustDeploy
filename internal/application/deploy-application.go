package application

import (
	"errors"
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

	eventsList := []adapter.EventServer{
		{
			Title:     "Build your application",
			EventType: "create_deploy",
		},
		{
			Title:     "Pull the traefik image",
			EventType: "create_deploy",
		},
		{
			Title:     "Run the traefik router",
			EventType: "create_deploy",
		},
		{
			Title:     "Start your application",
			EventType: "create_deploy",
		},
	}

	eventWrapper := adapter.EventDeployWrapper{
		DeployName:       deploy.Name,
		DeployId:         deploy.Id,
		EventsDeployList: eventsList,
		CurrentStep:      0,
	}

	// Build your application
	fmt.Println("Build your application")
	deployService.EventAdapter.SendNewDeployEvent(eventWrapper)
	fmt.Println("send: Build your application")

	err := deployService.DockerAdapter.BuildImage(deploy)
	if err != nil {
		deploy.Status = "Error"
		deployService.DatabaseAdapter.UpdateDeploy(*deploy)
		eventWrapper.SetStepError(err.Error())
		deployService.EventAdapter.SendNewDeployEvent(eventWrapper)
		return
	}

	// Pull the traefik image

	fmt.Println("Pull the traefik image")
	eventWrapper.NextStep()
	deployService.EventAdapter.SendNewDeployEvent(eventWrapper)
	fmt.Println("send: Pull the traefik image")
	err = deployService.DockerAdapter.PullTreafikImage()

	if err != nil {
		deploy.Status = "Error"
		deployService.DatabaseAdapter.UpdateDeploy(*deploy)
		eventWrapper.SetStepError(err.Error())
		deployService.EventAdapter.SendNewDeployEvent(eventWrapper)
		return
	}

	// Run the traefik router
	eventWrapper.NextStep()
	deployService.EventAdapter.SendNewDeployEvent(eventWrapper)
	err = deployService.DockerAdapter.RunRouter(deploy.Email)

	if err != nil {
		deploy.Status = "Error"
		deployService.DatabaseAdapter.UpdateDeploy(*deploy)
		eventWrapper.SetStepError(err.Error())
		deployService.EventAdapter.SendNewDeployEvent(eventWrapper)
		return
	}

	// Run your application
	eventWrapper.NextStep()
	deployService.EventAdapter.SendNewDeployEvent(eventWrapper)
	err = deployService.DockerAdapter.RunImage(deploy, appUrl)

	if err != nil {
		deploy.Status = "Error"
		deployService.DatabaseAdapter.UpdateDeploy(*deploy)
		eventWrapper.SetStepError(err.Error())
		deployService.EventAdapter.SendNewDeployEvent(eventWrapper)
	}

	eventWrapper.NextStep()
	deployService.EventAdapter.SendNewDeployEvent(eventWrapper)

	if deploy.EnableTls {
		deploy.Url = "https://" + appUrl
	} else {
		deploy.Url = "http://" + appUrl
	}

	fmt.Println(deploy.Url)
	deploy.Status = "Runing"
	deployService.DatabaseAdapter.UpdateDeploy(*deploy)
}

func DeployApplication(deployService *service.DeployService, newDeploy dto.NewDeployDto) (domain.Deploy, error) {
	server, err := deployService.DatabaseAdapter.GetServerById(newDeploy.ServerId)
	if err != nil {
		return domain.Deploy{}, err
	}

	if server.Domain == "" {
		return domain.Deploy{}, errors.New("server does not have domain")
	}

	pathToDir, err := filepath.Abs(newDeploy.PathToSource)

	if err != nil {
		return domain.Deploy{}, err
	}

	pathToDir = adapter.NewFilesystemAdapter().CleanPath(pathToDir)

	err = deployService.DockerAdapter.ConnectClient(server)

	if err != nil {
		return domain.Deploy{}, err
	}

	isFolder := adapter.NewFilesystemAdapter().IsFolder(pathToDir)
	DockerFileName := "Dockerfile"

	if !isFolder {
		DockerFileName = adapter.NewFilesystemAdapter().BaseDir(pathToDir)
		pathToDir = adapter.NewFilesystemAdapter().GetDir(pathToDir)
	}

	portEnv := make([]dto.Env, 1)

	portEnv[0] = dto.Env{
		Name:  "PORT",
		Value: "80",
	}

	newDeploy.Envs = append(portEnv, newDeploy.Envs...)

	Name := adapter.NewFilesystemAdapter().GetFolderName(pathToDir)

	SubDomain := ""

	deploys := deployService.DatabaseAdapter.GetDeployByServerId(server.Id)

	if len(deploys) > 0 {
		SubDomain = Name
	}

	deploy := domain.Deploy{
		Id:             utils.GenerateRandomPassword(5),
		Name:           Name,
		ServerId:       newDeploy.ServerId,
		PathToSource:   pathToDir,
		Status:         "Installing",
		EnableTls:      newDeploy.EnableTls,
		Email:          newDeploy.Email,
		Envs:           newDeploy.Envs,
		SubDomain:      SubDomain,
		DockerFileName: DockerFileName,
	}

	err = deployService.DatabaseAdapter.SaveDeploy(deploy)

	go runApplication(deployService, &deploy, server.Domain)

	return deploy, err
}
