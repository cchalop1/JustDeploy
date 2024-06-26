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

func (fs *FilesystemAdapter) GetCurrentPath() (string, error) {
	return os.Getwd()
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
	hooksFilePath := deploy.PathToSource + ".git/hooks/post-commit"
	fileContent := []byte("#!/bin/sh\njustdeploy --redeploy " + deploy.Id + "\n")

	err := os.WriteFile(hooksFilePath, fileContent, 0755)

	if err != nil {
		return err
	}

	fmt.Println("Create file ", deploy.PathToSource+".git/hooks/post-commit")
	return nil
}

func (fs *FilesystemAdapter) DeleteGitPostCommitHooks(deploy domain.Deploy) error {
	hooksFilePath := deploy.PathToSource + ".git/hooks/post-commit"

	err := os.Remove(hooksFilePath)

	if err != nil {
		return err
	}

	fmt.Println("Remove file", deploy.PathToSource+".git/hooks/post-commit")

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
