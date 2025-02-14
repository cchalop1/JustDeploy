package adapter

import (
	"fmt"
	"os/exec"
)

type GitAdapter struct{}

func NewGitAdapter() *GitAdapter {
	return &GitAdapter{}
}

func (g *GitAdapter) CloneRepository(repoUrl string, destPath string, token string) error {
	// If token is provided, inject it into the URL
	cloneUrl := repoUrl
	if token != "" {
		// Convert https://github.com/user/repo to https://<token>@github.com/user/repo
		cloneUrl = fmt.Sprintf("https://%s@github.com/%s", token, repoUrl)
	}

	fmt.Println("Cloning repository", cloneUrl, "to", destPath)

	cmd := exec.Command("git", "clone", cloneUrl, destPath)
	stdoutStderr, err := cmd.CombinedOutput()
	fmt.Printf("%s\n", stdoutStderr)
	return err
}
