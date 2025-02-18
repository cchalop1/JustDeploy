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
	"time"

	"cchalop1.com/deploy/internal/adapter/database"
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/domain"
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

func (d *DockerAdapter) ConnectClient(server domain.Server) error {
	var err error

	// if server.Ip != "localhost" {
	// 	serverCerts := server.GetCertsPath()

	// 	d.client, err = client.NewClientWithOpts(
	// 		client.WithHost("tcp://"+server.Ip+":2376"),
	// 		client.WithTLSClientConfig(
	// 			serverCerts.CaCertPath,
	// 			serverCerts.CertPath,
	// 			serverCerts.KeyPath,
	// 		),
	// 	)
	// } else {
	d.client, err = client.NewClientWithOpts(client.FromEnv)
	// }

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
func (d *DockerAdapter) BuildImage(service domain.Service) error {
	fmt.Println("Make a tar of", service.CurrentPath)
	tar, err := makeTar(service.CurrentPath)
	if err != nil {
		return err
	}
	fmt.Println("Tar created")
	fmt.Println("Building image", service.GetDockerName())

	buildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{service.GetDockerName()},
		Remove:     true,
	}

	buildResponse, err := d.client.ImageBuild(context.Background(), tar, buildOptions)
	if err != nil {
		fmt.Println("Error building image:", err)
		return err
	}
	defer buildResponse.Body.Close()

	scanner := bufio.NewScanner(buildResponse.Body)

	for scanner.Scan() {
		line := scanner.Text()
		var msg DockerMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			return fmt.Errorf("error decoding JSON: %v", err)
		}

		if msg.Stream != "" {
			fmt.Print(msg.Stream)
		}

		if msg.ErrorDetail.Message != "" || msg.Error != "" {
			errorMsg := msg.ErrorDetail.Message
			if errorMsg == "" {
				errorMsg = msg.Error
			}
			return fmt.Errorf("error building image: %s", errorMsg)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading build output: %v", err)
	}

	fmt.Println("Image built successfully")
	return nil

}

func (d *DockerAdapter) BuildNixpacksImage(service domain.Service) error {
	// Construct the nixpacks command

	nixpacksCmd := exec.Command("nixpacks", "build",
		"--name", service.GetDockerName(),
		service.CurrentPath)

	// Set up pipes for stdout and stderr
	stdout, _ := nixpacksCmd.StdoutPipe()
	stderr, _ := nixpacksCmd.StderrPipe()

	// Start the command
	if err := nixpacksCmd.Start(); err != nil {
		return fmt.Errorf("failed to start nixpacks build: %v", err)
	}

	// Create a scanner to read the output line by line
	scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	// Wait for the command to finish
	if err := nixpacksCmd.Wait(); err != nil {
		return fmt.Errorf("nixpacks build failed: %v", err)
	}

	fmt.Println("Nixpacks image built successfully")
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
	log.Printf("Pulling image: %s", image)
	reader, err := d.client.ImagePull(context.Background(), image, types.ImagePullOptions{})

	if err != nil {
		log.Printf("Error pulling image %s: %v", image, err)
		return err
	}

	_, err = io.Copy(io.Discard, reader)
	if err != nil {
		log.Printf("Error discarding image pull output for %s: %v", image, err)
		return err
	}

	log.Printf("Successfully pulled image: %s", image)
	return nil
}

func (d *DockerAdapter) RunRouter(email string) error {
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
			"--certificatesresolvers.myresolver.acme.email=" + email,
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
			// "/root/letsencrypt:/letsencrypt",
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

func envToSlice(envVars []dto.Env) []string {
	envSlice := make([]string, 0, len(envVars))
	for _, value := range envVars {
		if value.Name != "" && value.Value != "" {
			envSlice = append(envSlice, fmt.Sprintf("%s=%s", value.Name, value.Value))
		}
	}
	return envSlice
}

func (d *DockerAdapter) ConfigContainer(service domain.Service) container.Config {
	Image := service.ImageName

	if service.IsRepo {
		Image = service.GetDockerName()
	}

	config := container.Config{
		Image:    Image,
		Env:      envToSlice(service.Envs),
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

func (d *DockerAdapter) RunImage(service domain.Service, domain string) error {
	config := d.ConfigContainer(service)
	d.Stop(service.GetDockerName())
	d.Remove(service.GetDockerName())
	d.ExposeContainer(&config, ExposeContainerParams{
		IsTls:  false,
		Domain: domain,
		Port:   "3000",
	})
	con, err := d.client.ContainerCreate(context.Background(), &config, &container.HostConfig{}, &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"databases_default": {},
		},
	}, &v1.Platform{}, service.GetDockerName())

	if err != nil {
		fmt.Println(err)
		return err
	}

	d.client.ContainerStart(context.Background(), con.ID, container.StartOptions{})
	fmt.Printf("Container %s is started", con.ID)

	fmt.Println("Run image", service.GetDockerName())
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

func (d *DockerAdapter) GetLogsOfContainer(containerName string) ([]string, error) {
	logs, err := d.client.ContainerLogs(context.Background(), containerName, container.LogsOptions{ShowStdout: true, ShowStderr: true, Timestamps: true})
	if err != nil {
		return make([]string, 0), err
	}
	defer logs.Close()

	lines := make([]string, 0)
	buffer := make([]byte, 1024)

	for {
		n, err := logs.Read(buffer)
		if err != nil {
			break
		}
		clearLine := string(buffer[8:n])
		if n > 0 {
			lines = append(lines, clearLine)
		}
	}

	return lines, nil
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

func (d *DockerAdapter) RunServiceWithDeploy(service database.ServicesConfig, containerHostName string) {
	con, err := d.client.ContainerCreate(context.Background(), &service.Config, &container.HostConfig{},
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				"databases_default": {},
			},
		}, &v1.Platform{}, containerHostName)

	if err != nil {
		fmt.Println(err)
		return
	}

	d.client.ContainerStart(context.Background(), con.ID, container.StartOptions{})
	fmt.Printf("Container %s is started", con.ID)

	fmt.Println("Run image", service.Name)
}

func (d *DockerAdapter) GetLocalHostServer() domain.Server {
	return domain.Server{
		Id:       "local",
		Name:     "local",
		Ip:       "localhost",
		Domain:   "localhost",
		Password: nil,
		SshKey:   nil,
		// TODO: change the date
		CreatedDate: time.Now(),
		Status:      "active",
	}
}
