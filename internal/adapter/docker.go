package adapter

import (
	"context"
	"fmt"
	"io"
	"log"

	"cchalop1.com/deploy/internal"
	"cchalop1.com/deploy/internal/api/dto"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
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

func (d *DockerAdapter) ConnectClient(connectConfig dto.ConnectNewServerDto) error {
	caCertPath := internal.CERT_DOCKER_FOLDER + "/ca.pem"
	certPath := internal.CERT_DOCKER_FOLDER + "/cert.pem"
	keyPath := internal.CERT_DOCKER_FOLDER + "/key.pem"

	client, err := client.NewClientWithOpts(
		client.WithHost("tcp://"+connectConfig.Domain+":2376"),
		client.WithTLSClientConfig(caCertPath, certPath, keyPath),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}
	d.client = client
	fmt.Println("je suis connecter au docker")
	return nil
}

func makeTar(pathToDir string) (io.ReadCloser, error) {
	return archive.TarWithOptions(pathToDir, &archive.TarOptions{})
}

// BuildImage builds a Docker image from the specified Dockerfile and context directory.
func (d *DockerAdapter) BuildImage(imageName string, dockerfilePath string) {
	fmt.Println("Make a tar of", dockerfilePath)
	tar, err := makeTar(dockerfilePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	buildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{imageName},
		Remove:     true,
	}

	buildResponse, err := d.client.ImageBuild(context.Background(), tar, buildOptions)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer buildResponse.Body.Close()

	fmt.Println("Building image...")
	bytes, err := io.ReadAll(buildResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bytes))
	// check if sucesfull build image or not

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

func (d *DockerAdapter) PullTreafikImage() {
	treafikIsPulled, err := d.checkIsRouterImageIsPull()
	if treafikIsPulled || err != nil {
		return
	}

	fmt.Println("Pull image", ROUTER_NAME)
	reader, err := d.client.ImagePull(context.Background(), TRAEFIK_IMAGE, types.ImagePullOptions{})

	if err != nil {
		fmt.Println(err)
		return
	}

	io.Copy(io.Discard, reader)
}

func (d *DockerAdapter) RunRouter() {
	routerIsRuning, err := d.checkRouterIsRuning()
	if routerIsRuning || err != nil {
		return
	}
	email := "clement.chalopin@gmail.com"

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
			// "8080/tcp": struct{}{},
		},
	}

	portMap := nat.PortMap{
		"80/tcp":  []nat.PortBinding{{HostIP: "", HostPort: "80"}},
		"443/tcp": []nat.PortBinding{{HostIP: "", HostPort: "443"}},
		// "8080/tcp": []nat.PortBinding{{HostIP: "", HostPort: "8080"}},
	}

	con, err := d.client.ContainerCreate(context.Background(), &config, &container.HostConfig{
		Binds: []string{
			"/root/letsencrypt:/letsencrypt",
			"/var/run/docker.sock:/var/run/docker.sock:ro",
		},
		NetworkMode:  "default",
		PortBindings: portMap,
	}, nil, &v1.Platform{}, ROUTER_NAME)

	if err != nil {
		fmt.Println(err)
		return
	}

	d.client.ContainerStart(context.Background(), con.ID, types.ContainerStartOptions{})
	fmt.Printf("Container %s is started", con.ID)

	fmt.Println("Run image", ROUTER_NAME)

}

func envToSlice(envVars []dto.Env) []string {
	envSlice := make([]string, 0, len(envVars))
	for _, value := range envVars {
		if value.Name != "" && value.Secret != "" {
			envSlice = append(envSlice, fmt.Sprintf("%s=%s", value.Name, value.Secret))
		}
	}
	return envSlice
}

func (d *DockerAdapter) RunImage(deployConfig dto.DeployConfigDto) {
	Name := deployConfig.AppConfig.Name

	Labels := map[string]string{
		"traefik.enable":                                "true",
		"traefik.http.routers." + Name + ".rule":        "Host(`" + deployConfig.ServerConfig.Domain + "`)",
		"traefik.http.routers." + Name + ".entrypoints": "web",
	}

	if deployConfig.AppConfig.EnableTls {
		Labels["traefik.http.routers."+Name+".tls"] = "true"
		Labels["traefik.http.routers."+Name+".tls.certresolver"] = "myresolver"
		Labels["traefik.http.routers."+Name+".entrypoints"] = "websecure"
	}

	config := container.Config{
		Image:  Name,
		Labels: Labels,
		Env:    envToSlice(deployConfig.AppConfig.Envs),
	}

	con, err := d.client.ContainerCreate(context.Background(), &config, &container.HostConfig{}, nil, &v1.Platform{}, Name)

	if err != nil {
		fmt.Println(err)
		return
	}

	d.client.ContainerStart(context.Background(), con.ID, types.ContainerStartOptions{})
	fmt.Printf("Container %s is started", con.ID)

	fmt.Println("Run image", Name)
}

func (d *DockerAdapter) Delete(appName string, stopRouter bool) {
	d.Stop(appName)
	d.Remove(appName)
	if stopRouter {
		d.Stop(ROUTER_NAME)
		d.Remove(ROUTER_NAME)
	}
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

func (d *DockerAdapter) GetLogsOfContainer(containerName string) []string {
	logs, err := d.client.ContainerLogs(context.Background(), containerName, container.LogsOptions{ShowStdout: true})
	if err != nil {
		fmt.Println("Failed to read container logs: ", err)
		return make([]string, 0)
	}
	defer logs.Close()

	lines := make([]string, 0)
	buf := make([]byte, 1024)

	for {
		n, err := logs.Read(buf)
		if n > 0 {
			lines = append(lines, string(buf[:n]))
		}
		if err != nil {
			break
		}
	}

	return lines
}
