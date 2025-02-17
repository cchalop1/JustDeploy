package application

import (
	"fmt"
	"path/filepath"

	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
)

func runApplication(deployService *service.DeployService, deploy *domain.Deploy, domain string) {

	// // eventsList := []adapter.EventServer{
	// // 	{
	// // 		Title:     "Build your application",
	// // 		EventType: "create_deploy",
	// // 	},
	// // 	{
	// // 		Title:     "Pull the traefik image",
	// // 		EventType: "create_deploy",
	// // 	},
	// // 	{
	// // 		Title:     "Run the traefik router",
	// // 		EventType: "create_deploy",
	// // 	},
	// // 	{
	// // 		Title:     "Start your application",
	// // 		EventType: "create_deploy",
	// // 	},
	// // }

	// // eventWrapper := adapter.EventDeployWrapper{
	// // 	DeployName:       deploy.Name,
	// // 	DeployId:         deploy.Id,
	// // 	EventsDeployList: eventsList,
	// // 	CurrentStep:      0,
	// // }

	// // Build your application
	// fmt.Println("Build your application")
	// // go deployService.EventAdapter.SendNewDeployEvent(eventWrapper)
	// fmt.Println("send: Build your application")

	// var err error

	// isDockerfile := deployService.FilesystemAdapter.FindDockerFile(deploy.PathToSource)
	// if isDockerfile {
	// 	err := deployService.DockerAdapter.BuildImage(deploy)
	// 	if err != nil {
	// 		deploy.Status = "Error"
	// 		deployService.DatabaseAdapter.UpdateDeploy(*deploy)
	// 		// eventWrapper.SetStepError(err.Error())
	// 		// deployService.EventAdapter.SendNewDeployEvent(eventWrapper)
	// 		return
	// 	}
	// } else {
	// 	server, err := deployService.DatabaseAdapter.GetServerById(deploy.ServerId)

	// 	if err != nil {
	// 		deploy.Status = "Error"
	// 		deployService.DatabaseAdapter.UpdateDeploy(*deploy)
	// 		// eventWrapper.SetStepError(err.Error())
	// 		// deployService.EventAdapter.SendNewDeployEvent(eventWrapper)
	// 		return
	// 	}

	// 	err = deployService.DockerAdapter.BuildNixpacksImage(deploy, server)
	// 	if err != nil {
	// 		deploy.Status = "Error"
	// 		deployService.DatabaseAdapter.UpdateDeploy(*deploy)
	// 		// eventWrapper.SetStepError(err.Error())
	// 		// deployService.EventAdapter.SendNewDeployEvent(eventWrapper)
	// 		return
	// 	}
	// }
	// // Pull the traefik image

	// fmt.Println("Pull the traefik image")
	// // eventWrapper.NextStep()
	// // go deployService.EventAdapter.SendNewDeployEvent(eventWrapper)
	// fmt.Println("send: Pull the traefik image")
	// err = deployService.DockerAdapter.PullTreafikImage()

	// if err != nil {
	// 	deploy.Status = "Error"
	// 	deployService.DatabaseAdapter.UpdateDeploy(*deploy)
	// 	// eventWrapper.SetStepError(err.Error())
	// 	// deployService.EventAdapter.SendNewDeployEvent(eventWrapper)
	// 	return
	// }

	// // Run the traefik router
	// // eventWrapper.NextStep()
	// // go deployService.EventAdapter.SendNewDeployEvent(eventWrapper)
	// err = deployService.DockerAdapter.RunRouter(deploy.Email)

	// if err != nil {
	// 	deploy.Status = "Error"
	// 	deployService.DatabaseAdapter.UpdateDeploy(*deploy)
	// 	// eventWrapper.SetStepError(err.Error())
	// 	// deployService.EventAdapter.SendNewDeployEvent(eventWrapper)
	// 	return
	// }

	// // Run your application
	// // eventWrapper.NextStep()
	// // go deployService.EventAdapter.SendNewDeployEvent(eventWrapper)
	// err = deployService.DockerAdapter.RunImage(deploy, appUrl)

	// if err != nil {
	// 	deploy.Status = "Error"
	// 	deployService.DatabaseAdapter.UpdateDeploy(*deploy)
	// 	// eventWrapper.SetStepError(err.Error())
	// 	// deployService.EventAdapter.SendNewDeployEvent(eventWrapper)
	// }

	// // eventWrapper.NextStep()
	// // go deployService.EventAdapter.SendNewDeployEvent(eventWrapper)

	// if deploy.EnableTls {
	// 	deploy.Url = "https://" + appUrl
	// } else {
	// 	deploy.Url = "http://" + appUrl
	// }

	// fmt.Println(deploy.Url)
	// deploy.Status = "Runing"
	// deployService.DatabaseAdapter.UpdateDeploy(*deploy)
}

func deployOneService(deployService *service.DeployService, serviceToDeploy domain.Service) error {
	pathToDir, err := filepath.Abs(serviceToDeploy.CurrentPath)

	if err != nil {
		return err
	}

	pathToDir = adapter.NewFilesystemAdapter().CleanPath(pathToDir)

	isFolder := adapter.NewFilesystemAdapter().IsFolder(pathToDir)
	DockerFileName := "Dockerfile"

	if !isFolder {
		fmt.Println("Is not a folder")
		DockerFileName = adapter.NewFilesystemAdapter().BaseDir(pathToDir)
		pathToDir = adapter.NewFilesystemAdapter().GetDir(pathToDir)
	}

	portEnv := make([]dto.Env, 1)

	portEnv[0] = dto.Env{
		Name:  "PORT",
		Value: "80",
	}

	fmt.Println("Path to dir: ", pathToDir)
	fmt.Println("Docker file name: ", DockerFileName)

	serviceToDeploy.Envs = append(portEnv, serviceToDeploy.Envs...)

	isDockerfile := deployService.FilesystemAdapter.FindDockerFile(serviceToDeploy.CurrentPath)
	if isDockerfile {
		err := deployService.DockerAdapter.BuildImage(serviceToDeploy)
		if err != nil {
			return fmt.Errorf("error building Docker image: %w", err)
		}
	} else {
		// TODO: nixpacks build
	}

	domain := "localhost"

	err = deployService.DockerAdapter.RunImage(serviceToDeploy, domain)

	if err != nil {
		return fmt.Errorf("error running Docker image: %w", err)
	}

	// if deploy.EnableTls {
	// 	deploy.Url = "https://" + appUrl
	// } else {
	// 	deploy.Url = "http://" + appUrl
	// }

	return nil
}

func DeployApplication(deployService *service.DeployService) error {
	services := deployService.DatabaseAdapter.GetServices()

	server := deployService.DatabaseAdapter.GetServer()

	err := deployService.DockerAdapter.ConnectClient(server)

	if err != nil {
		return fmt.Errorf("error connecting Docker client: %w", err)
	}

	err = deployService.DockerAdapter.PullTreafikImage()

	if err != nil {
		return fmt.Errorf("error pulling Traefik image: %w", err)
	}

	err = deployService.DockerAdapter.RunRouter("clement.chalopin@gmail.com")

	if err != nil {
		return fmt.Errorf("error running Traefik router: %w", err)
	}

	// if server.Domain == "" {
	// 	return errors.New("Server does not have domain")
	// }

	for _, service := range services {
		err = deployOneService(deployService, service)
		if err != nil {
			return fmt.Errorf("error deploying service: %w", err)
		}
	}

	return nil
}
