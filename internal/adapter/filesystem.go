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
	"cchalop1.com/deploy/internal/adapter/database"
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/domain"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
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

type serviceConfig struct {
	Image       string            `yaml:"image"`
	Ports       []string          `yaml:"ports"`
	Environment map[string]string `yaml:"environment"`
	Volumes     []string          `yaml:"volumes"`
	Name        string            `yaml:"container_name"`
	Cmd         []string          `yaml:"command"`
}

type composeConfig struct {
	Services map[string]serviceConfig `yaml:"services"`
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

func filterComposeServiceToArray(services map[string]serviceConfig) []database.ServicesConfig {
	servicesArray := []database.ServicesConfig{}

	for key, value := range services {
		envs := []dto.Env{}

		for key := range value.Environment {
			envs = append(envs, dto.Env{Name: key, Value: "", IsSecret: true})
		}

		ports := nat.PortSet{}

		for _, port := range value.Ports {
			ports[nat.Port(port)] = struct{}{}
		}

		// volumes := []string{}

		// for _, volume := range value.Volumes {
		// 	volumes = append(volumes, strings.Split(volume, ":")[1])
		// }

		servicesArray = append(servicesArray, database.ServicesConfig{
			Name: key,
			Icon: "compose",
			Env:  envs,
			Config: container.Config{
				Image:        value.Image,
				Cmd:          value.Cmd,
				ExposedPorts: ports,
			},
		})
	}

	return servicesArray
}

func (fs *FilesystemAdapter) GetComposeConfigOfDeploy(pathToSource string) ([]database.ServicesConfig, error) {
	// TODO: try all the compose file name like (docker-compose.yml, docker-compose.yaml, compose.yml, compose.yaml)
	cfg, err := parsComposeFile(pathToSource + "/docker-compose.yml")

	if err != nil {
		return nil, err
	}

	if len(cfg.Services) == 0 {
		return nil, errors.New("no services found in the docker-compose.yml")
	}

	services := filterComposeServiceToArray(cfg.Services)

	return services, nil
}

// .env file management
func (fs *FilesystemAdapter) GenerateDotEnvFile(project *domain.Project) error {
	// Full path to the .env file
	for _, service := range project.Services {
		if service.IsDevContainer {

			envFilePath := service.CurrentPath + "/.env"

			// Check if .env file exists
			_, err := os.Stat(envFilePath)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					// File does not exist, create it
					fmt.Println(".env file does not exist, creating a new one")
				} else {
					// Other error when trying to check file
					return err
				}
			}

			// Open the file in append mode if it exists, otherwise create it
			file, err := os.OpenFile(envFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			defer file.Close()

			// Write environment variables to the file
			for _, env := range service.Envs {
				_, err := file.WriteString(env.Name + "=" + env.Value + "\n")
				if err != nil {
					return err
				}
			}

		}
	}

	return nil
}

func (fs *FilesystemAdapter) RemoveEnvsFromDotEnvFile(project *domain.Project, envToRemove []dto.Env) error {
	for _, service := range project.Services {
		if service.IsDevContainer {

			envFilePath := service.CurrentPath + "/.env"

			// Check if .env file exists
			_, err := os.Stat(envFilePath)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					// File does not exist, nothing to remove
					return nil
				}
				// Other error when trying to check file
				return err
			}

			// Read the entire file
			content, err := os.ReadFile(envFilePath)
			if err != nil {
				return err
			}

			// Create a map of environment variables to remove for quick lookup
			envsToRemove := make(map[string]bool)
			for _, env := range envToRemove {
				envsToRemove[env.Name] = true
			}

			// Process the file line by line
			var newLines []string
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				trimmedLine := strings.TrimSpace(line)
				if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
					// Keep empty lines and comments
					newLines = append(newLines, line)
					continue
				}

				parts := strings.SplitN(trimmedLine, "=", 2)
				if len(parts) < 2 {
					// Keep lines that don't look like environment variable declarations
					newLines = append(newLines, line)
					continue
				}

				key := strings.TrimSpace(parts[0])
				if !envsToRemove[key] {
					// Keep lines for environment variables that are not in the removal list
					newLines = append(newLines, line)
				}
			}

			// Write the updated content back to the file
			newContent := strings.Join(newLines, "\n")
			err = os.WriteFile(envFilePath, []byte(newContent), 0644)
			if err != nil {
				return err
			}

		}
	}

	return nil
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
