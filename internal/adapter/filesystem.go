package adapter

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"cchalop1.com/deploy/internal"
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/internal/utils"
	"gopkg.in/yaml.v3"
)

type FilesystemAdapter struct {
}

func NewFilesystemAdapter() *FilesystemAdapter {
	return &FilesystemAdapter{}
}

func (fs *FilesystemAdapter) GetFolderName(path string) string {
	cleanedPath := filepath.Clean(path)

	projectName := filepath.Base(cleanedPath)

	return projectName
}

func (fs *FilesystemAdapter) GetFullPathToProject(path string) string {
	cleanedPath := filepath.Clean(path)

	fullPathToProject, err := filepath.Abs(cleanedPath)

	if err != nil {
		panic("Error to find the path of the project")

	}
	return fullPathToProject
}

func (fs *FilesystemAdapter) CleanPath(path string) string {
	if path[len(path)-1] != '/' {
		return path + "/"
	}
	return path
}

func (fs *FilesystemAdapter) GetCurrentPath() string {
	path, err := os.Getwd()
	if err != nil {
		return ""
	}
	return path
}

func (fs *FilesystemAdapter) FindDockerFile(pathToFolder string) bool {
	entries, err := os.ReadDir(pathToFolder)
	if err != nil {
		panic("Error to read the directory you have")
	}

	for _, e := range entries {
		if e.Name() == "Dockerfile" || strings.HasPrefix(e.Name(), "Dockerfile.") {
			return true
		}
	}
	return false
}

func (fs *FilesystemAdapter) FindDockerComposeFile(pathToFolder string) bool {
	entries, err := os.ReadDir(pathToFolder)
	if err != nil {
		panic("Error to read the directory you have")
	}

	for _, e := range entries {
		if e.Name() == "docker-compose.yml" || e.Name() == "compose.yml" {
			return true
		}
	}
	return false
}

func (fs *FilesystemAdapter) LoadEnvsFromFileSystem(pathToFolder string) []dto.Env {
	file, err := os.Open(pathToFolder + "/.env")
	var envs []dto.Env = []dto.Env{}

	if err != nil {
		return envs
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			env := dto.Env{
				Name:  strings.TrimSpace(parts[0]),
				Value: strings.TrimSpace(parts[1]),
			}
			envs = append(envs, env)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error while scanning the .env file")
	}

	return envs
}

func (fs *FilesystemAdapter) CopyFileToRemoteServer(sourcePath string, serverIp string) error {
	cmd := exec.Command("scp", sourcePath, "root@"+serverIp+":/root/")
	fmt.Println(cmd.String())
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(string(stdout))
	return nil
}

func (fs *FilesystemAdapter) CreateGitPostCommitHooks(deploy domain.Deploy) error {
	// hooksFilePath := deploy.PathToSource + ".git/hooks/post-commit"
	// fileContent := []byte("#!/bin/sh\njustdeploy --redeploy " + deploy.Id + "\n")

	// err := os.WriteFile(hooksFilePath, fileContent, 0755)

	// if err != nil {
	// 	return err
	// }

	// fmt.Println("Create file ", deploy.PathToSource+".git/hooks/post-commit")
	return nil
}

func (fs *FilesystemAdapter) DeleteGitPostCommitHooks(deploy domain.Deploy) error {
	// hooksFilePath := deploy.PathToSource + ".git/hooks/post-commit"

	// err := os.Remove(hooksFilePath)

	// if err != nil {
	// 	return err
	// }

	// fmt.Println("Remove file", deploy.PathToSource+".git/hooks/post-commit")

	return nil
}

func (fs *FilesystemAdapter) RemoveDockerCertificateByServerId(serverId string) error {
	pathLocalCertDir := internal.CERT_DOCKER_FOLDER + "/" + serverId + "/"
	return os.Remove(pathLocalCertDir)
}

func (fs *FilesystemAdapter) IsFolder(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

func (fs *FilesystemAdapter) BaseDir(path string) string {
	return filepath.Base(path)
}

func (fs *FilesystemAdapter) GetDir(path string) string {
	arr := strings.Split(path, "/")
	if len(arr) > 1 {
		arr = arr[:len(arr)-2] // remove the last element
		return strings.Join(arr, "/") + "/"
	}
	return "/"
}

// Get docker compose config from the file

type buildConfig struct {
	Context    string `yaml:"context"`
	Dockerfile string `yaml:"dockerfile"`
}

type ComposeServiceConfig struct {
	Image       string      `yaml:"image"`
	Ports       []string    `yaml:"ports"`
	Environment interface{} `yaml:"environment"`
	Volumes     []string    `yaml:"volumes"`
	Name        string      `yaml:"container_name"`
	Cmd         []string    `yaml:"command"`
	Build       interface{} `yaml:"build"`
	DependsOn   []string    `yaml:"depends_on"`
}

type composeConfig struct {
	Services map[string]ComposeServiceConfig `yaml:"services"`
}

func parsComposeFile(pathToComposeFile string) (*composeConfig, error) {
	file, err := os.ReadFile(pathToComposeFile)

	if err != nil {
		return nil, err
	}

	var cfg composeConfig

	yaml.Unmarshal(file, &cfg)

	return &cfg, nil
}

func (s *ComposeServiceConfig) HasBuild() bool {
	return s.Build != nil
}

// GetEnvironmentVariables converts the Environment field to a slice of dto.Env
func (s *ComposeServiceConfig) GetEnvironmentVariables() []dto.Env {
	var envVars []dto.Env

	switch env := s.Environment.(type) {
	case map[string]interface{}:
		// Handle map format: environment: { VAR1: value1, VAR2: value2 }
		for key, value := range env {
			strValue := ""
			if value != nil {
				strValue = fmt.Sprintf("%v", value)
			}
			envVars = append(envVars, dto.Env{
				Name:  key,
				Value: strValue,
			})
		}
	case map[interface{}]interface{}:
		// Handle map format with interface keys: environment: { VAR1: value1, VAR2: value2 }
		for key, value := range env {
			strKey := fmt.Sprintf("%v", key)
			strValue := ""
			if value != nil {
				strValue = fmt.Sprintf("%v", value)
			}
			envVars = append(envVars, dto.Env{
				Name:  strKey,
				Value: strValue,
			})
		}
	case []interface{}:
		// Handle array format: environment: [ VAR1=value1, VAR2=value2 ] or [ VAR1, VAR2 ]
		for _, item := range env {
			if strItem, ok := item.(string); ok {
				// Check if it's in the format VAR=value
				parts := strings.SplitN(strItem, "=", 2)
				if len(parts) == 2 {
					envVars = append(envVars, dto.Env{
						Name:  parts[0],
						Value: parts[1],
					})
				} else {
					// It's just a variable name without a value
					envVars = append(envVars, dto.Env{
						Name:  parts[0],
						Value: "",
					})
				}
			}
		}
	case []string:
		// Handle string array format: environment: [ VAR1=value1, VAR2=value2 ] or [ VAR1, VAR2 ]
		for _, item := range env {
			parts := strings.SplitN(item, "=", 2)
			if len(parts) == 2 {
				envVars = append(envVars, dto.Env{
					Name:  parts[0],
					Value: parts[1],
				})
			} else {
				envVars = append(envVars, dto.Env{
					Name:  parts[0],
					Value: "",
				})
			}
		}
	}

	return envVars
}

func (fs *FilesystemAdapter) GetComposeConfigOfDeploy(pathToSource string) (map[string]ComposeServiceConfig, error) {
	// TODO: try all the compose file name like (docker-compose.yml, docker-compose.yaml, compose.yml, compose.yaml)
	cfg, err := parsComposeFile(pathToSource + "/docker-compose.yml")

	if err != nil {
		return nil, err
	}

	if len(cfg.Services) == 0 {
		return nil, errors.New("no services found in the docker-compose.yml")
	}

	return cfg.Services, nil
}

// Get list of folders
func (fs *FilesystemAdapter) GetFolders(path string) ([]dto.PathDto, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("path does not exist: %s", path)
	}

	return fs.getFoldersRecursive(path, 0, 2)
}

func (fs *FilesystemAdapter) getFoldersRecursive(path string, currentDepth, maxDepth int) ([]dto.PathDto, error) {
	if currentDepth >= maxDepth {
		return nil, nil
	}

	var result []dto.PathDto

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			fullPath := filepath.Join(path, file.Name())
			if strings.HasPrefix(file.Name(), ".") {
				continue
			}

			folder := dto.PathDto{
				Name:     file.Name(),
				FullPath: fullPath,
			}

			// subFolders, err := fs.getFoldersRecursive(fullPath, currentDepth+1, maxDepth)
			// if err != nil {
			// 	return nil, err
			// }

			// folder.Folders = subFolders
			result = append(result, folder)
		}
	}

	return result, nil
}

func (f *FilesystemAdapter) GetTempDir() string {
	return os.TempDir() + "/deploy-" + utils.GenerateRandomPassword(5)
}
