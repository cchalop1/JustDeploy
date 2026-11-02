package adapter

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
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

func (g *GitAdapter) CloneRepository(repo, dest, token string) error {
	return g.CloneRepositoryBranch(repo, dest, token, "")
}
func (g *GitAdapter) CloneRepositoryBranch(repo, dest, token, branch string) error {
	if !regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-]*/[A-Za-z0-9_.-]+$`).MatchString(repo) || strings.Contains(repo, "..") || dest == "" {
		return fmt.Errorf("invalid repository")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	args := []string{"clone", "--depth", "1"}
	if branch != "" {
		if err := exec.CommandContext(ctx, "git", "check-ref-format", "--branch", branch).Run(); err != nil {
			return fmt.Errorf("invalid branch")
		}
		args = append(args, "--branch", branch)
	}
	args = append(args, "--", "https://github.com/"+repo+".git", dest)
	cmd := exec.CommandContext(ctx, "git", args...)
	env := []string{}
	for _, v := range os.Environ() {
		if !strings.HasPrefix(v, "GIT_") {
			env = append(env, v)
		}
	}
	cmd.Env = append(env, "GIT_TERMINAL_PROMPT=0")
	if token != "" {
		header := "AUTHORIZATION: basic " + base64.StdEncoding.EncodeToString([]byte("x-access-token:"+token))
		cmd.Env = append(cmd.Env, "GIT_CONFIG_COUNT=1", "GIT_CONFIG_KEY_0=http.https://github.com/.extraheader", "GIT_CONFIG_VALUE_0="+header)
	}
	// Git error output may contain credentials supplied by a remote; return a bounded, secret-free error.
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}
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

func (g *GitAdapter) GetBranch(path string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	output, err := exec.CommandContext(ctx, "git", "-C", path, "branch", "--show-current").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
