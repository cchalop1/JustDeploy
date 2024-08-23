package adapter

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"

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
const ROUTER_NAME = "treafik"

func (d *DockerAdapter) ConnectClient(server domain.Server) error {
	serverCerts := server.GetCertsPath()

	client, err := client.NewClientWithOpts(
		client.WithHost("tcp://"+server.Ip+":2376"),
		client.WithTLSClientConfig(
			serverCerts.CaCertPath,
			serverCerts.CertPath,
			serverCerts.KeyPath,
		),
		// client.WithTimeout(3*time.Second),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = client.ContainerList(context.Background(), container.ListOptions{})

	if err != nil {
		fmt.Println(err)
		return err
	}

	d.client = client
	fmt.Println("I'm connected to docker of the server ", server.Name, " With domain ", server.Ip)
	return nil
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
func (d *DockerAdapter) BuildImage(deploy *domain.Deploy) error {
	fmt.Println("Make a tar of", deploy.PathToSource)
	tar, err := makeTar(deploy.PathToSource)
	if err != nil {
		return err
	}

	buildOptions := types.ImageBuildOptions{
		Dockerfile: deploy.DockerFileName,
		Tags:       []string{deploy.GetDockerName()},
		Remove:     true,
	}

	buildResponse, err := d.client.ImageBuild(context.Background(), tar, buildOptions)
	if err != nil {
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
	if treafikIsPulled || err != nil {
		return err
	}

	d.PullImage(TRAEFIK_IMAGE)
	return nil
}

func (d *DockerAdapter) PullImage(image string) {
	fmt.Println("Pull image", image)
	reader, err := d.client.ImagePull(context.Background(), image, types.ImagePullOptions{})

	if err != nil {
		fmt.Println(err)
		return
	}

	io.Copy(io.Discard, reader)
}

func (d *DockerAdapter) RunRouter(email string) error {
	routerIsRuning, err := d.checkRouterIsRuning()
	if routerIsRuning || err != nil {
		return nil
	}

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
			"/root/letsencrypt:/letsencrypt",
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

func (d *DockerAdapter) RunImage(deploy *domain.Deploy, domain string) error {
	Name := deploy.GetDockerName()

	Labels := map[string]string{
		"traefik.enable":                                "true",
		"traefik.http.routers." + Name + ".rule":        "Host(`" + domain + "`)",
		"traefik.http.routers." + Name + ".entrypoints": "web",
	}

	if deploy.EnableTls {
		Labels["traefik.http.routers."+Name+".tls"] = "true"
		Labels["traefik.http.routers."+Name+".tls.certresolver"] = "myresolver"
		Labels["traefik.http.routers."+Name+".entrypoints"] = "websecure"
	}

	config := container.Config{
		Image:  Name,
		Labels: Labels,
		Env:    envToSlice(deploy.Envs),
	}

	con, err := d.client.ContainerCreate(context.Background(), &config, &container.HostConfig{}, &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"databases_default": {},
		},
	}, &v1.Platform{}, Name)

	if err != nil {
		return err
	}

	d.client.ContainerStart(context.Background(), con.ID, container.StartOptions{})
	fmt.Printf("Container %s is started", con.ID)

	fmt.Println("Run image", Name)
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

func (d *DockerAdapter) RunService(service database.ServicesConfig, containerHostName string) {
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

	d.client.ContainerStart(context.Background(), con.ID, types.ContainerStartOptions{})
	fmt.Printf("Container %s is started", con.ID)

	fmt.Println("Run image", service.Name)
}
