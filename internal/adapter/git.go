package adapter

import (
	"fmt"
	"os"
	"os/exec"
)

type GitAdapter struct{}

func NewGitAdapter() *GitAdapter {
	return &GitAdapter{}
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
