package adapter

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/internal/utils"
	"github.com/containerd/errdefs"
	"github.com/moby/go-archive"
	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

type DockerAdapter struct{ client *client.Client }

const TRAEFIK_IMAGE = "traefik:v3.6.1"
const ROUTER_NAME = "traefik"
const proxyNetwork = "justdeploy-proxy"

func NewDockerAdapter() *DockerAdapter { return &DockerAdapter{} }
func (d *DockerAdapter) ConnectClient() error {
	if d.client != nil {
		return nil
	}
	var err error
	d.client, err = client.New(client.FromEnv, client.WithTimeout(35*time.Minute))
	return err
}

func makeTar(path string) (io.ReadCloser, error) {
	excludes := []string{".git", "**/.git", ".env", ".env.*", "**/.env", "**/.env.*", "**/*.pem", "**/*.key"}
	if data, err := os.ReadFile(filepath.Join(path, ".dockerignore")); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				excludes = append(excludes, line)
			}
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	// Mandatory exclusions come last, so negated dockerignore rules cannot re-include secrets.
	excludes = append(excludes, ".git", "**/.git", ".env", ".env.*", "**/.env", "**/.env.*", "**/*.pem", "**/*.key")
	return archive.TarWithOptions(path, &archive.TarOptions{ExcludePatterns: excludes})
}

type DockerMessage struct {
	Stream      string `json:"stream"`
	ErrorDetail struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errorDetail"`
	Error string `json:"error"`
}

func consumeDockerStream(reader io.Reader, logs *domain.LogCollector) error {
	decoder := json.NewDecoder(reader)
	for {
		var msg DockerMessage
		err := decoder.Decode(&msg)
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("invalid Docker response: %w", err)
		}
		if msg.ErrorDetail.Message != "" {
			return errors.New(msg.ErrorDetail.Message)
		}
		if msg.Error != "" {
			return errors.New(msg.Error)
		}
		if logs != nil && msg.Stream != "" {
			logs.AddLog(msg.Stream)
		}
	}
}
func (d *DockerAdapter) BuildImage(s domain.Service, logs *domain.LogCollector) error {
	tar, err := makeTar(s.CurrentPath)
	if err != nil {
		return err
	}
	defer tar.Close()
	dockerfile := s.Dockerfile
	if dockerfile == "" {
		dockerfile, err = FindDockerfile(s.CurrentPath)
		if err != nil {
			return err
		}
	}
	if dockerfile == "" {
		return errors.New("Dockerfile not found")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	response, err := d.client.ImageBuild(ctx, tar, client.ImageBuildOptions{Dockerfile: dockerfile, Tags: []string{s.GetDockerName()}, Remove: true})
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return consumeDockerStream(response.Body, logs)
}

type buildLogWriter struct {
	mu      sync.Mutex
	logs    *domain.LogCollector
	pending string
}

func (w *buildLogWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.pending += string(p)
	for {
		i := strings.IndexByte(w.pending, '\n')
		if i < 0 {
			break
		}
		w.logs.AddLog(w.pending[:i])
		w.pending = w.pending[i+1:]
	}
	if len(w.pending) > 64*1024 {
		w.logs.AddLog(w.pending[:64*1024])
		w.pending = ""
	}
	return len(p), nil
}
func (d *DockerAdapter) BuildNixpacksImage(s domain.Service, logs *domain.LogCollector) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	// Sanitize the workspace as a build context too: Nixpacks must not receive .git credentials or .env files.
	tar, err := makeTar(s.CurrentPath)
	if err != nil {
		return err
	}
	defer tar.Close()
	dir, err := os.MkdirTemp("", "justdeploy-build-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	if err := archive.Untar(tar, dir, &archive.TarOptions{}); err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, "nixpacks", "build", "--name", s.GetDockerName(), dir)
	writer := &buildLogWriter{logs: logs}
	cmd.Stdout = writer
	cmd.Stderr = writer
	err = cmd.Run()
	if writer.pending != "" {
		logs.AddLog(writer.pending)
	}
	if err != nil {
		return fmt.Errorf("nixpacks build failed: %w", err)
	}
	return nil
}
func (d *DockerAdapter) PullImage(image string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	reader, err := d.client.ImagePull(ctx, image, client.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()
	return consumeDockerStream(reader, nil)
}
func (d *DockerAdapter) EnsureImage(image string) error {
	_, err := d.client.ImageInspect(context.Background(), image)
	if err == nil {
		return nil
	}
	if !errdefs.IsNotFound(err) {
		return err
	}
	return d.PullImage(image)
}
func (d *DockerAdapter) PullTreafikImage() error             { return d.EnsureImage(TRAEFIK_IMAGE) }
func (d *DockerAdapter) CheckRouterIsRunning() (bool, error) { return d.IsServiceRunning(ROUTER_NAME) }
func (d *DockerAdapter) IsServiceRunning(name string) (bool, error) {
	info, err := d.inspect(context.Background(), name)
	if errdefs.IsNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.State != nil && info.State.Running, nil
}
func (d *DockerAdapter) EnsureNetwork(name string, internal bool) error {
	_, err := d.client.NetworkInspect(context.Background(), name, client.NetworkInspectOptions{})
	if err == nil {
		return nil
	}
	if !errdefs.IsNotFound(err) {
		return err
	}
	_, err = d.client.NetworkCreate(context.Background(), name, client.NetworkCreateOptions{Driver: "bridge", Internal: internal, Labels: map[string]string{"justdeploy.managed": "true"}})
	return err
}
func (d *DockerAdapter) RunRouterWithServer(server domain.Server) error {
	if err := d.EnsureNetwork(proxyNetwork, false); err != nil {
		return err
	}
	fingerprint := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%s|%t|%s|%s", TRAEFIK_IMAGE, server.UseHttps, server.Email, os.Getenv("JUSTDEPLOY_DOCKER_PROXY")))))
	old, err := d.inspect(context.Background(), ROUTER_NAME)
	exists := err == nil
	if err != nil && !errdefs.IsNotFound(err) {
		return err
	}
	if exists && old.State != nil && old.State.Running && old.Config != nil && old.Config.Labels["justdeploy.router-config"] == fingerprint {
		return nil
	}
	cmd := []string{"--providers.docker=true", "--providers.docker.exposedbydefault=false", "--providers.docker.network=" + proxyNetwork, "--entrypoints.web.address=:80", "--entrypoints.websecure.address=:443"}
	if server.UseHttps {
		if server.Email == "" {
			return errors.New("ACME email is required for HTTPS")
		}
		cmd = append(cmd, "--certificatesresolvers.myresolver.acme.tlschallenge=true", "--certificatesresolvers.myresolver.acme.email="+server.Email, "--certificatesresolvers.myresolver.acme.storage=/letsencrypt/acme.json", "--entrypoints.web.http.redirections.entrypoint.to=websecure", "--entrypoints.web.http.redirections.entrypoint.scheme=https")
	}
	// A restricted Docker socket proxy can be supplied instead of mounting the raw socket.
	binds := []string{"letsencrypt:/letsencrypt"}
	if endpoint := os.Getenv("JUSTDEPLOY_DOCKER_PROXY"); endpoint != "" {
		cmd = append(cmd, "--providers.docker.endpoint="+endpoint)
	} else {
		binds = append(binds, "/var/run/docker.sock:/var/run/docker.sock:ro")
	}
	config := &container.Config{Image: TRAEFIK_IMAGE, Cmd: cmd, ExposedPorts: network.PortSet{network.MustParsePort("80/tcp"): {}, network.MustParsePort("443/tcp"): {}}, Labels: map[string]string{"justdeploy.managed": "true", "justdeploy.router-config": fingerprint}}
	host := &container.HostConfig{Binds: binds, PortBindings: network.PortMap{network.MustParsePort("80/tcp"): {{HostPort: "80"}}, network.MustParsePort("443/tcp"): {{HostPort: "443"}}}, RestartPolicy: container.RestartPolicy{Name: "unless-stopped"}, Resources: container.Resources{Memory: 256 * 1024 * 1024, PidsLimit: int64ptr(200)}}
	wasRunning := exists && old.State != nil && old.State.Running
	if exists {
		if err := d.Stop(ROUTER_NAME); err != nil {
			return err
		}
		if _, err := d.client.ContainerRename(context.Background(), old.ID, client.ContainerRenameOptions{NewName: ROUTER_NAME + "-previous"}); err != nil {
			if wasRunning {
				_ = d.Start(old.ID)
			}
			return err
		}
	}
	restore := func(cause error) error {
		if err := d.Delete(ROUTER_NAME); err != nil {
			return errors.Join(cause, err)
		}
		if exists {
			if _, err := d.client.ContainerRename(context.Background(), old.ID, client.ContainerRenameOptions{NewName: ROUTER_NAME}); err != nil {
				return errors.Join(cause, err)
			}
			if wasRunning {
				if err := d.Start(old.ID); err != nil {
					return errors.Join(cause, err)
				}
			}
		}
		return cause
	}
	con, err := d.client.ContainerCreate(context.Background(), client.ContainerCreateOptions{Config: config, HostConfig: host, NetworkingConfig: &network.NetworkingConfig{EndpointsConfig: map[string]*network.EndpointSettings{proxyNetwork: {}}}, Name: ROUTER_NAME})
	if err != nil {
		return restore(err)
	}
	if err := d.Start(con.ID); err != nil {
		return restore(err)
	}
	if err := d.waitReady(con.ID); err != nil {
		return restore(err)
	}
	if exists {
		return d.Remove(old.ID)
	}
	return nil
}

func int64ptr(n int64) *int64 { return &n }
func ProjectNetwork(s domain.Service) string {
	id := s.ProjectID
	if id == "" {
		id = s.Id
	}
	return "justdeploy-project-" + strings.ToLower(id)
}
func (d *DockerAdapter) ConfigContainer(s domain.Service) container.Config {
	image := s.ImageName
	if s.Type == "github_repo" {
		image = s.GetDockerName()
	}
	return container.Config{Image: image, Env: utils.EnvToSlice(s.Envs), Hostname: s.GetDockerName(), Cmd: s.Cmd, Labels: map[string]string{"justdeploy.managed": "true", "justdeploy.service": s.Id, "traefik.enable": "false"}}
}

type ExposeContainerParams struct {
	IsTls        bool
	Domain, Port string
}

func (d *DockerAdapter) ExposeContainer(config *container.Config, p ExposeContainerParams) {
	// Router names must be unique even when several services use the same image.
	name := config.Hostname
	config.Labels["traefik.enable"] = "true"
	config.Labels["traefik.docker.network"] = proxyNetwork
	prefix := "traefik.http.routers." + name
	config.Labels[prefix+".rule"] = "Host(`" + p.Domain + "`)"
	config.Labels[prefix+".entrypoints"] = "web"
	config.Labels[prefix+".service"] = name
	config.Labels["traefik.http.services."+name+".loadbalancer.server.port"] = p.Port
	if p.IsTls {
		config.Labels[prefix+".tls"] = "true"
		config.Labels[prefix+".tls.certresolver"] = "myresolver"
		config.Labels[prefix+".entrypoints"] = "websecure"
	}
}
func DefaultDataTargets(image string) []string {
	base := strings.Split(strings.Split(image, "@")[0], ":")[0]
	base = filepath.Base(base)
	switch base {
	case "postgres":
		if strings.HasPrefix(image, "postgres:18") {
			return []string{"/var/lib/postgresql"}
		}
		return []string{"/var/lib/postgresql/data"}
	case "mongo":
		return []string{"/data/db", "/data/configdb"}
	case "redis", "minio":
		return []string{"/data"}
	case "ollama":
		return []string{"/root/.ollama"}
	}
	return nil
}
func serviceMounts(s domain.Service, old container.InspectResponse) ([]mount.Mount, error) {
	mounts := []mount.Mount{}
	// Adopt all existing data mounts, including anonymous Docker volumes, rather than hiding them with empty named volumes.
	for _, m := range old.Mounts {
		if m.Type != mount.TypeVolume && m.Type != mount.TypeBind {
			return nil, fmt.Errorf("cannot safely preserve mount %s", m.Destination)
		}
		source := m.Source
		if m.Type == mount.TypeVolume {
			source = m.Name
		}
		mounts = append(mounts, mount.Mount{Type: m.Type, Source: source, Target: m.Destination, ReadOnly: !m.RW})
	}
	targets := DefaultDataTargets(s.ImageName)
	for _, v := range s.Volumes {
		targets = append(targets, v.Target)
	}
	for _, target := range targets {
		exists := false
		for _, m := range mounts {
			if m.Target == target || strings.HasPrefix(target, m.Target+"/") {
				exists = true
			}
		}
		if exists {
			continue
		}
		if old.ID != "" && s.Type != "github_repo" {
			return nil, fmt.Errorf("existing data at %s has no reusable mount; back up and migrate before replacing this service", target)
		}
		name := "justdeploy-" + strings.ToLower(s.Id) + "-" + strings.ReplaceAll(strings.Trim(target, "/"), "/", "-")
		for _, v := range s.Volumes {
			if v.Target == target && v.Name != "" {
				name = v.Name
			}
		}
		mounts = append(mounts, mount.Mount{Type: mount.TypeVolume, Source: name, Target: target})
	}
	if s.Type == "database" && len(mounts) == 0 {
		return nil, errors.New("database storage target is unknown; configure a named data volume before deploying")
	}
	return mounts, nil
}

// RunImageWithTLS retains the previous container until the replacement has passed its startup check.
// Application images are built before this method; stateful services reuse existing mounts.
func (d *DockerAdapter) RunImageWithTLS(s domain.Service, domainName string, useHTTPS bool) error {
	config := d.ConfigContainer(s)
	project := ProjectNetwork(s)
	if err := d.EnsureNetwork(project, false); err != nil {
		return err
	}
	endpoints := map[string]*network.EndpointSettings{project: {Aliases: []string{s.Name, s.GetDockerName()}}}
	if s.ExposeSettings.IsExposed && domainName != "" {
		if err := d.EnsureNetwork(proxyNetwork, false); err != nil {
			return err
		}
		endpoints[proxyNetwork] = &network.EndpointSettings{}
		d.ExposeContainer(&config, ExposeContainerParams{useHTTPS, domainName, s.ExposeSettings.ExposePort})
	}
	old, err := d.inspect(context.Background(), s.GetDockerName())
	if err != nil && !errdefs.IsNotFound(err) {
		return err
	}
	exists := err == nil
	mounts, err := serviceMounts(s, old)
	if err != nil {
		return err
	}
	for _, m := range mounts {
		if m.Type == mount.TypeVolume {
			if _, err := d.client.VolumeCreate(context.Background(), client.VolumeCreateOptions{Name: m.Source, Labels: map[string]string{"justdeploy.service": s.Id}}); err != nil {
				return err
			}
		}
	}
	host := &container.HostConfig{Mounts: mounts, RestartPolicy: container.RestartPolicy{Name: "unless-stopped"}, Resources: container.Resources{Memory: 1024 * 1024 * 1024, NanoCPUs: 2_000_000_000, PidsLimit: int64ptr(512)}, SecurityOpt: []string{"no-new-privileges:true"}}
	backup := s.GetDockerName() + "-previous"
	wasRunning := exists && old.State != nil && old.State.Running
	if exists {
		if err := d.Stop(s.GetDockerName()); err != nil {
			return err
		}
		if _, err := d.client.ContainerRename(context.Background(), old.ID, client.ContainerRenameOptions{NewName: backup}); err != nil {
			if wasRunning {
				_ = d.Start(old.ID)
			}
			return err
		}
	}
	restore := func(cause error) error {
		if err := d.Delete(s.GetDockerName()); err != nil {
			return errors.Join(cause, fmt.Errorf("remove failed replacement: %w", err))
		}
		if exists {
			if _, err := d.client.ContainerRename(context.Background(), old.ID, client.ContainerRenameOptions{NewName: s.GetDockerName()}); err != nil {
				return errors.Join(cause, err)
			}
			if wasRunning {
				if err := d.Start(old.ID); err != nil {
					return errors.Join(cause, fmt.Errorf("rollback failed: %w", err))
				}
			}
		}
		return cause
	}
	con, err := d.client.ContainerCreate(context.Background(), client.ContainerCreateOptions{Config: &config, HostConfig: host, NetworkingConfig: &network.NetworkingConfig{EndpointsConfig: endpoints}, Name: s.GetDockerName()})
	if err != nil {
		return restore(err)
	}
	if err := d.Start(con.ID); err != nil {
		return restore(err)
	}
	if err := d.waitReady(con.ID); err != nil {
		return restore(err)
	}
	if exists {
		if err := d.Remove(old.ID); err != nil {
			return fmt.Errorf("new service is running but previous container cleanup failed: %w", err)
		}
	}
	return nil
}
func (d *DockerAdapter) waitReady(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	started := time.Now()
	for {
		info, err := d.inspect(ctx, id)
		if err != nil {
			return err
		}
		if info.State == nil || !info.State.Running || info.State.Restarting || info.RestartCount > 0 {
			return errors.New("container exited during startup")
		}
		if info.State.Health != nil {
			if info.State.Health.Status == "healthy" {
				return nil
			}
			if info.State.Health.Status == "unhealthy" {
				return errors.New("container health check failed")
			}
		} else if time.Since(started) >= 2*time.Second {
			return nil
		}
		select {
		case <-ctx.Done():
			return errors.New("container did not become healthy in time")
		case <-time.After(250 * time.Millisecond):
		}
	}
}
func (d *DockerAdapter) Stop(name string) error {
	_, err := d.client.ContainerStop(context.Background(), name, client.ContainerStopOptions{})
	if errdefs.IsNotFound(err) {
		return nil
	}
	return err
}
func (d *DockerAdapter) Remove(name string) error {
	// Never remove data volumes as a side effect of a container operation.
	_, err := d.client.ContainerRemove(context.Background(), name, client.ContainerRemoveOptions{RemoveVolumes: false})
	if errdefs.IsNotFound(err) {
		return nil
	}
	return err
}
func (d *DockerAdapter) Start(name string) error {
	_, err := d.client.ContainerStart(context.Background(), name, client.ContainerStartOptions{})
	return err
}
func (d *DockerAdapter) Delete(name string) error {
	if err := d.Stop(name); err != nil {
		return err
	}
	return d.Remove(name)
}
func (d *DockerAdapter) GetLogsOfContainer(name string) ([]domain.Logs, error) {
	info, err := d.inspect(context.Background(), name)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	logs, err := d.client.ContainerLogs(ctx, name, client.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Timestamps: true, Tail: "1000"})
	if err != nil {
		return nil, err
	}
	defer logs.Close()
	var out limitedLogBuffer
	if info.Config != nil && info.Config.Tty {
		_, err = io.Copy(&out, logs)
	} else {
		_, err = stdcopy.StdCopy(&out, &out, logs)
	}
	if err != nil {
		return nil, err
	}
	result := []domain.Logs{}
	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Buffer(make([]byte, 4096), 1024*1024)
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), " ", 2)
		if len(parts) != 2 {
			continue
		}
		date, err := time.Parse(time.RFC3339Nano, parts[0])
		if err != nil {
			continue
		}
		result = append(result, domain.Logs{Date: date.Format(time.RFC3339), Message: parts[1]})
	}
	return result, scanner.Err()
}

type limitedLogBuffer struct{ bytes.Buffer }

func (b *limitedLogBuffer) Write(p []byte) (int, error) {
	if b.Len()+len(p) > 4*1024*1024 {
		return 0, errors.New("log response exceeds 4 MiB")
	}
	return b.Buffer.Write(p)
}

func (d *DockerAdapter) inspect(ctx context.Context, id string) (container.InspectResponse, error) {
	result, err := d.client.ContainerInspect(ctx, id, client.ContainerInspectOptions{})
	return result.Container, err
}
