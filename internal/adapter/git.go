package adapter

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type GitAdapter struct{}

func NewGitAdapter() *GitAdapter {
	return &GitAdapter{}
}

type CommitInfo struct {
	Hash    string `json:"hash"`
	Message string `json:"message"`
	Author  string `json:"author"`
	Date    string `json:"date"`
}

func (g *GitAdapter) CloneRepository(repoUrl string, destPath string, token string) error {
	if repoUrl == "" || destPath == "" {
		return fmt.Errorf("repository URL and destination path cannot be empty")
	}

	cloneUrl := fmt.Sprintf("https://github.com/%s", repoUrl)
	if token != "" {
		cloneUrl = fmt.Sprintf("https://oauth2:%s@github.com/%s", token, repoUrl)
	}

	fmt.Println("Cloning repository", cloneUrl, "to", destPath)

	cmd := exec.Command("git", "clone", cloneUrl, destPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	fmt.Println("Cloning successful!")
	return nil
}

// GetLastCommitInfo récupère les informations du dernier commit depuis un répertoire Git local
func (g *GitAdapter) GetLastCommitInfo(repoPath string) (*CommitInfo, error) {
	if repoPath == "" {
		return nil, fmt.Errorf("repository path cannot be empty")
	}

	// Vérifier que le répertoire existe et contient un dépôt Git
	if _, err := os.Stat(repoPath + "/.git"); os.IsNotExist(err) {
		return nil, fmt.Errorf("not a git repository: %s", repoPath)
	}

	// Récupérer le hash du commit
	hashCmd := exec.Command("git", "-C", repoPath, "rev-parse", "HEAD")
	hashOutput, err := hashCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get commit hash: %w", err)
	}
	hash := strings.TrimSpace(string(hashOutput))

	// Récupérer le message du commit
	messageCmd := exec.Command("git", "-C", repoPath, "log", "-1", "--pretty=format:%s")
	messageOutput, err := messageCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get commit message: %w", err)
	}
	message := strings.TrimSpace(string(messageOutput))

	// Récupérer l'auteur du commit
	authorCmd := exec.Command("git", "-C", repoPath, "log", "-1", "--pretty=format:%an")
	authorOutput, err := authorCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get commit author: %w", err)
	}
	author := strings.TrimSpace(string(authorOutput))

	// Récupérer la date du commit
	dateCmd := exec.Command("git", "-C", repoPath, "log", "-1", "--pretty=format:%ci")
	dateOutput, err := dateCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get commit date: %w", err)
	}
	date := strings.TrimSpace(string(dateOutput))

	return &CommitInfo{
		Hash:    hash,
		Message: message,
		Author:  author,
		Date:    date,
	}, nil
}

// GetCommitGithubUrl génère l'URL GitHub pour un commit spécifique
func (g *GitAdapter) GetCommitGithubUrl(repoFullName string, commitHash string) string {
	return fmt.Sprintf("https://github.com/%s/commit/%s", repoFullName, commitHash)
}
