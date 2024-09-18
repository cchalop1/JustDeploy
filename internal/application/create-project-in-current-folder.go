package application

import (
	"fmt"
	"os"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
)

func CreateProjectCurrentFolder(deployService *service.DeployService) (string, error) {
	currentPath := deployService.FilesystemAdapter.GetCurrentPath()

	project := deployService.DatabaseAdapter.GetProjectByPath(currentPath)

	Name := deployService.FilesystemAdapter.GetFolderName(currentPath)

	if project != nil {
		fmt.Printf("Projet déjà existant: %s\n", project.Name)
		return project.Id, nil
	}

	createProjectDto := dto.CreateProjectDto{
		Name: Name,
		Path: currentPath,
	}

	projectId, err := CreateProject(deployService, createProjectDto)
	if err != nil {
		fmt.Println("Erreur lors de la création du projet:", err)
		os.Exit(1)
		return "", err
	}

	// createServiceDto := dto.CreateServiceDto{
	// 	ServiceName: Name,
	// 	ProjectId:   &projectId,
	// 	LocalPath:   &currentPath,
	// }

	// CreateService(deployService, createServiceDto)
	fmt.Printf("Projet créé avec succès avec l'ID: %s\n", projectId)
	return projectId, nil
}
