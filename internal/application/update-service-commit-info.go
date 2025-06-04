package application

import (
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
)

// UpdateServiceCommitInfo met à jour les informations du dernier commit pour un service GitHub
func UpdateServiceCommitInfo(deployService *service.DeployService, serviceToUpdate *domain.Service) error {
	// Vérifier que c'est un service GitHub
	if serviceToUpdate.Type != "github_repo" {
		return nil // Ne rien faire pour les services non-GitHub
	}

	// Récupérer les informations du commit depuis le répertoire local
	commitInfo, err := deployService.GitAdapter.GetLastCommitInfo(serviceToUpdate.CurrentPath)
	if err != nil {
		// Si on ne peut pas récupérer les informations de commit, on continue sans erreur
		// car cela ne doit pas empêcher le déploiement
		return nil
	}

	// Mettre à jour les informations de commit dans le service
	serviceToUpdate.LastCommit = domain.CommitInfo{
		Hash:    commitInfo.Hash,
		Message: commitInfo.Message,
		Author:  commitInfo.Author,
		Date:    commitInfo.Date,
	}

	// Sauvegarder le service avec les nouvelles informations de commit
	return deployService.DatabaseAdapter.SaveService(*serviceToUpdate)
}
