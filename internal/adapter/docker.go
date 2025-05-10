package adapter

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"cchalop1.com/deploy/internal/adapter/database"
	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/internal/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type DockerAdapter struct {
	client *client.Client
}

func NewDockerAdapter() *DockerAdapter {
	return &DockerAdapter{}
}

const TRAEFIK_IMAGE = "traefik"
const ROUTER_NAME = "traefik"

func (d *DockerAdapter) ConnectClient() error {
	var err error

	d.client, err = client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		fmt.Println("Error creating Docker client:", err)
	}

	return err
}

func makeTar(pathToDir string) (io.ReadCloser, error) {
	return archive.TarWithOptions(pathToDir, &archive.TarOptions{})
}

type DockerMessage struct {
	Stream      string `json:"stream"`
	ErrorDetail struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errorDetail"`
	Error string `json:"error"`
}

// BuildImage builds a Docker image from the specified Dockerfile and context directory.
func (d *DockerAdapter) BuildImage(service domain.Service, logCollector *domain.LogCollector) error {
	logCollector.AddLog(fmt.Sprintf("Making tar of %s", service.CurrentPath))
	tar, err := makeTar(service.CurrentPath)
	if err != nil {
		logCollector.AddLog(fmt.Sprintf("Error creating tar: %v", err))
		return err
	}
	logCollector.AddLog("Tar created")
	logCollector.AddLog(fmt.Sprintf("Building image %s", service.GetDockerName()))

	buildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{service.GetDockerName()},
		Remove:     true,
	}

	buildResponse, err := d.client.ImageBuild(context.Background(), tar, buildOptions)
	if err != nil {
		logCollector.AddLog(fmt.Sprintf("Error building image: %v", err))
		return err
	}
	defer buildResponse.Body.Close()

	scanner := bufio.NewScanner(buildResponse.Body)

	for scanner.Scan() {
		line := scanner.Text()
		var msg DockerMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			logCollector.AddLog(fmt.Sprintf("Error decoding JSON: %v", err))
			return fmt.Errorf("error decoding JSON: %v", err)
		}

		if msg.Stream != "" {
			logCollector.AddLog(msg.Stream)
		}

		if msg.ErrorDetail.Message != "" || msg.Error != "" {
			errorMsg := msg.ErrorDetail.Message
			if errorMsg == "" {
				errorMsg = msg.Error
			}
			logCollector.AddLog(fmt.Sprintf("Error building image: %s", errorMsg))
			return fmt.Errorf("error building image: %s", errorMsg)
		}
	}

	if err := scanner.Err(); err != nil {
		logCollector.AddLog(fmt.Sprintf("Error reading build output: %v", err))
		return fmt.Errorf("error reading build output: %v", err)
	}

	logCollector.AddLog("Image built successfully")
	return nil
}

func (d *DockerAdapter) BuildNixpacksImage(service domain.Service, logCollector *domain.LogCollector) error {
	logCollector.AddLog(fmt.Sprintf("Starting Nixpacks build for %s", service.GetDockerName()))

	nixpacksCmd := exec.Command("nixpacks", "build",
		"--name", service.GetDockerName(),
		service.CurrentPath)

	stdout, _ := nixpacksCmd.StdoutPipe()
	stderr, _ := nixpacksCmd.StderrPipe()

	if err := nixpacksCmd.Start(); err != nil {
		logCollector.AddLog(fmt.Sprintf("Failed to start nixpacks build: %v", err))
		return fmt.Errorf("failed to start nixpacks build: %v", err)
	}

	scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
	for scanner.Scan() {
		logCollector.AddLog(scanner.Text())
	}

	if err := nixpacksCmd.Wait(); err != nil {
		logCollector.AddLog(fmt.Sprintf("Nixpacks build failed: %v", err))
		return fmt.Errorf("nixpacks build failed: %v", err)
	}

	logCollector.AddLog("Nixpacks image built successfully")
	return nil
}

func (d *DockerAdapter) checkIsRouterImageIsPull() (bool, error) {
	imageList, err := d.client.ImageList(context.Background(), types.ImageListOptions{})

	if err != nil {

		return false, err
	}

	for _, image := range imageList {
		if len(image.RepoTags) > 0 && image.RepoTags[0] == TRAEFIK_IMAGE {
			fmt.Println("Image", TRAEFIK_IMAGE, "already exists")
			return true, nil
		}
	}
	return false, nil
}

func (d *DockerAdapter) checkRouterIsRuning() (bool, error) {
	containerList, err := d.client.ContainerList(context.Background(), types.ContainerListOptions{})

	if err != nil {
		return false, err
	}

	for _, container := range containerList {
		if container.Names[0] == "/treafik" {
			fmt.Println("Router", ROUTER_NAME, "is already runing")
			return true, nil
		}
	}
	return false, nil
}

// CheckRouterIsRunning checks if the Traefik router is running
func (d *DockerAdapter) CheckRouterIsRunning() (bool, error) {
	return d.checkRouterIsRuning()
}

// IsServiceRunning checks if a service with the given name is running
func (d *DockerAdapter) IsServiceRunning(serviceName string) (bool, error) {
	containerList, err := d.client.ContainerList(context.Background(), types.ContainerListOptions{})

	if err != nil {
		return false, err
	}

	for _, container := range containerList {
		for _, name := range container.Names {
			// Remove leading slash from container name
			if name == "/"+serviceName || name[1:] == serviceName {
				return true, nil
			}
		}
	}
	return false, nil
}

func (d *DockerAdapter) PullTreafikImage() error {
	treafikIsPulled, err := d.checkIsRouterImageIsPull()
	if err != nil {
		return err
	}

	if treafikIsPulled {
		return nil
	}

	d.PullImage(TRAEFIK_IMAGE)
	return nil
}

func (d *DockerAdapter) PullImage(image string) error {
	d.ConnectClient()
	log.Printf("Pulling image: %s", image)
	reader, err := d.client.ImagePull(context.Background(), image, types.ImagePullOptions{})

	if err != nil {
		log.Printf("Error pulling image %s: %v", image, err)
		return err
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		var msg DockerMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			log.Printf("Error decoding JSON: %v", err)
			log.Println(line) // Print raw line if JSON parsing fails
			continue
		}

		if msg.Stream != "" {
			fmt.Print(msg.Stream)
		} else if msg.Error != "" {
			log.Printf("Error: %s", msg.Error)
		} else {
			// Print the raw line if no stream or error field
			log.Println(line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading image pull output for %s: %v", image, err)
		return err
	}

	log.Printf("Successfully pulled image: %s", image)
	return nil
}

// RunRouterWithServer runs the Traefik router with server configuration
func (d *DockerAdapter) RunRouterWithServer(server domain.Server) error {
	d.Stop(ROUTER_NAME)
	d.Remove(ROUTER_NAME)
	d.client.NetworkCreate(context.Background(), "databases_default", types.NetworkCreate{})

	config := container.Config{
		Image: TRAEFIK_IMAGE,
		Cmd: []string{
			"--api.insecure=true",
			"--providers.docker=true",
			"--providers.docker.exposedbydefault=false",
			"--entrypoints.web.address=:80",
			"--entrypoints.websecure.address=:443",
			"--certificatesresolvers.myresolver.acme.tlschallenge=true",
			"--certificatesresolvers.myresolver.acme.email=" + server.Email,
			"--certificatesresolvers.myresolver.acme.storage=/letsencrypt/acme.json",
		},
		ExposedPorts: nat.PortSet{
			"80/tcp":  struct{}{},
			"443/tcp": struct{}{},
		},
	}

	portMap := nat.PortMap{
		"80/tcp":  []nat.PortBinding{{HostIP: "", HostPort: "80"}},
		"443/tcp": []nat.PortBinding{{HostIP: "", HostPort: "443"}},
	}

	con, err := d.client.ContainerCreate(context.Background(), &config, &container.HostConfig{
		Binds: []string{
			"letsencrypt:/letsencrypt",
			"/var/run/docker.sock:/var/run/docker.sock:ro",
		},
		NetworkMode:  "default",
		PortBindings: portMap,
	}, &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"databases_default": {},
		},
	}, &v1.Platform{}, ROUTER_NAME)

	if err != nil {
		return err
	}

	d.client.ContainerStart(context.Background(), con.ID, container.StartOptions{})
	fmt.Printf("Container %s is started", con.ID)

	fmt.Println("Run image", ROUTER_NAME)
	return nil
}

func (d *DockerAdapter) ConfigContainer(service domain.Service) container.Config {
	Image := service.ImageName

	if service.Type == "github_repo" {
		Image = service.GetDockerName()
	}

	config := container.Config{
		Image:    Image,
		Env:      utils.EnvToSlice(service.Envs),
		Hostname: service.GetDockerName(),
	}
	return config
}

type ExposeContainerParams struct {
	IsTls  bool
	Domain string
	Port   string
}

func (d *DockerAdapter) ExposeContainer(containersConfig *container.Config, exposeContainerParams ExposeContainerParams) {

	Labels := map[string]string{
		"traefik.enable": "true",
		"traefik.http.routers." + containersConfig.Image + ".rule":                      "Host(`" + exposeContainerParams.Domain + "`)",
		"traefik.http.routers." + containersConfig.Image + ".entrypoints":               "web",
		"traefik.http.services." + containersConfig.Image + ".loadbalancer.server.port": exposeContainerParams.Port,
	}

	if exposeContainerParams.IsTls {
		Labels["traefik.http.routers."+containersConfig.Image+".tls"] = "true"
		Labels["traefik.http.routers."+containersConfig.Image+".tls.certresolver"] = "myresolver"
		Labels["traefik.http.routers."+containersConfig.Image+".entrypoints"] = "websecure"
	}

	containersConfig.Labels = Labels
}

// RunImageWithTLS runs a container with optional TLS settings from the server
func (d *DockerAdapter) RunImageWithTLS(service domain.Service, domain string, useHttps bool) error {
	config := d.ConfigContainer(service)
	d.Stop(service.GetDockerName())
	d.Remove(service.GetDockerName())
	if domain != "" {
		d.ExposeContainer(&config, ExposeContainerParams{
			IsTls:  useHttps,
			Domain: domain,
			Port:   service.ExposeSettings.ExposePort,
		})
	}
	con, err := d.client.ContainerCreate(context.Background(), &config, &container.HostConfig{}, &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"databases_default": {},
		},
	}, &v1.Platform{}, service.GetDockerName())

	if err != nil {
		fmt.Println(err)
		return err
	}

	err = d.client.ContainerStart(context.Background(), con.ID, container.StartOptions{})
	if err != nil {
		return fmt.Errorf("failed to start container: %v", err)
	}

	fmt.Printf("Container %s is started\n", con.ID)
	fmt.Println("Run image", service.GetDockerName())

	// Vérifier si le conteneur est bien en cours d'exécution
	// Attendre un court instant pour laisser le temps au conteneur de démarrer
	time.Sleep(2 * time.Second)

	// containerInfo, err := d.client.ContainerInspect(context.Background(), con.ID)
	// if err != nil {
	// 	return fmt.Errorf("failed to inspect container: %v", err)
	// }

	// if !containerInfo.State.Running {
	// 	logs, logErr := d.GetLogsOfContainer(service.GetDockerName())
	// 	if logErr == nil && len(logs) > 0 {
	// 		return fmt.Errorf("container failed to start. Logs: %s", strings.Join(logs, "\n"))
	// 	}
	// 	return fmt.Errorf("container failed to start. Exit code: %d, Error: %s",
	// 		containerInfo.State.ExitCode, containerInfo.State.Error)
	// }

	return nil
}

func (d *DockerAdapter) Delete(appName string) {
	d.Stop(appName)
	d.Remove(appName)
}

func (d *DockerAdapter) Stop(appName string) {
	d.client.ContainerStop(context.Background(), appName, container.StopOptions{})
}

func (d *DockerAdapter) Remove(appName string) {
	d.client.ContainerRemove(context.Background(), appName, container.RemoveOptions{})
}

func (d *DockerAdapter) Start(appName string) {
	d.client.ContainerStart(context.Background(), appName, container.StartOptions{})
}

func (d *DockerAdapter) GetLogsOfContainer(containerName string) ([]domain.Logs, error) {
	logs, err := d.client.ContainerLogs(context.Background(), containerName, container.LogsOptions{ShowStdout: true, ShowStderr: true, Timestamps: true})
	if err != nil {
		return []domain.Logs{}, err
	}
	defer logs.Close()

	parsedLogs := make([]domain.Logs, 0)
	scanner := bufio.NewScanner(logs)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 8 {
			continue
		}

		// Skip the first 8 bytes (Docker log header)
		content := line[8:]

		// Parse the timestamp and message
		// Format: 2025-05-10T03:25:42.310893047Z message
		parts := strings.SplitN(content, " ", 2)
		if len(parts) != 2 {
			continue
		}

		timestamp := parts[0]
		message := parts[1]

		// Parse the timestamp
		date, err := time.Parse(time.RFC3339Nano, timestamp)
		if err != nil {
			continue
		}

		parsedLogs = append(parsedLogs, domain.Logs{
			Date:    date.Format(time.RFC3339),
			Message: message,
		})
	}

	if err := scanner.Err(); err != nil {
		return parsedLogs, err
	}

	return parsedLogs, nil
}

func (d *DockerAdapter) RunService(service database.ServicesConfig, exposedPort string, containerHostName string) {
	internalPort := strconv.Itoa(service.DefaultPort) + "/tcp"
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			nat.Port(internalPort): []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: exposedPort,
				},
			},
		},
	}

	con, err := d.client.ContainerCreate(context.Background(), &service.Config, hostConfig,
		&network.NetworkingConfig{}, &v1.Platform{}, containerHostName)

	if err != nil {
		fmt.Println(err)
		return
	}

	d.client.ContainerStart(context.Background(), con.ID, container.StartOptions{})
	fmt.Printf("Container %s is started", con.ID)

	fmt.Println("Run image", service.Name)
}
