package application

import (
	"fmt"
	"os"

	"cchalop1.com/deploy/internal/api/http/dto"
	"cchalop1.com/deploy/internal/api/service"
)

func CreateProjectForCurrentFolder(deployService *service.DeployService) (string, error) {
	currentPath := deployService.FilesystemAdapter.GetCurrentPath()
	project := deployService.DatabaseAdapter.GetProjectByPath(currentPath)

	if project == nil {
		createProjectDto := dto.CreateProjectDto{
			Name: deployService.FilesystemAdapter.GetFolderName(currentPath),
			Path: currentPath,
		}

		projectId, err := CreateProject(deployService, createProjectDto)
		if err != nil {
			fmt.Println("Erreur lors de la création du projet:", err)
			os.Exit(1)
		}

		CreateApp(deployService, dto.CreateAppDto{Path: currentPath, ProjectId: projectId})
		fmt.Printf("Projet créé avec succès avec l'ID: %s\n", projectId)
		return projectId, nil
	} else {
		fmt.Printf("Projet déjà existant: %s\n", project.Name)
		return project.Id, nil
	}
}
