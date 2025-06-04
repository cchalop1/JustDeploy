package application

import (
	"cchalop1.com/deploy/internal/api/service"
)

type ServiceCommitInfo struct {
	Hash      string `json:"hash"`
	Message   string `json:"message"`
	Author    string `json:"author"`
	Date      string `json:"date"`
	GithubUrl string `json:"githubUrl"`
}

// GetServiceCommitInfo récupère les informations du dernier commit déployé pour un service GitHub
func GetServiceCommitInfo(deployService *service.DeployService, serviceId string) (*ServiceCommitInfo, error) {
	// Récupérer le service par son ID
	service, err := deployService.DatabaseAdapter.GetServiceById(serviceId)
	if err != nil {
		return nil, err
	}

	// Vérifier que c'est un service GitHub
	if service.Type != "github_repo" {
		return nil, nil // Retourner nil pour les services non-GitHub sans erreur
	}

	// Vérifier si nous avons des informations de commit stockées
	if service.LastCommit.Hash == "" {
		return nil, nil // Pas encore de commit déployé
	}

	// Générer l'URL GitHub pour le commit
	githubUrl := deployService.GitAdapter.GetCommitGithubUrl(service.FullName, service.LastCommit.Hash)

	// Afficher seulement les 8 premiers caractères du hash
	shortHash := service.LastCommit.Hash
	if len(shortHash) > 8 {
		shortHash = shortHash[:8]
	}

	return &ServiceCommitInfo{
		Hash:      shortHash,
		Message:   service.LastCommit.Message,
		Author:    service.LastCommit.Author,
		Date:      service.LastCommit.Date,
		GithubUrl: githubUrl,
	}, nil
}
