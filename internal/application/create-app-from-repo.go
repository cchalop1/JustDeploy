package application

import (
	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"cchalop1.com/deploy/internal/utils"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func CreateServiceFromGithubRepo(ds *service.DeployService, repo string) (domain.Service, error) {
	if !ValidRepo(repo) {
		return domain.Service{}, errors.New("repository must have the form owner/name")
	}
	if !ds.Operation.TryLock() {
		return domain.Service{}, domain.ErrConflict
	}
	defer ds.Operation.Unlock()
	settings, err := ds.DatabaseAdapter.ReadSettings()
	if err != nil {
		return domain.Service{}, err
	}
	if settings.GithubToken == "" {
		return domain.Service{}, errors.New("GitHub integration is not configured")
	}
	token := settings.GithubToken
	if !ds.GithubAdapter.IsTokenValid(token) {
		token, err = renewGithubToken(ds, settings.GithubAppSettings)
		if err != nil {
			return domain.Service{}, err
		}
	}
	dir, err := ds.FilesystemAdapter.CreateWorkspace()
	if err != nil {
		return domain.Service{}, err
	}
	success := false
	defer func() {
		if !success {
			_ = os.RemoveAll(dir)
		}
	}()
	if err := ds.GitAdapter.CloneRepository(repo, dir, token); err != nil {
		return domain.Service{}, err
	}
	name := strings.ToLower(strings.Split(repo, "/")[1])
	name = strings.ReplaceAll(name, ".", "-")
	id := utils.GenerateRandomPassword(16)
	app := domain.Service{Id: id, Name: name, FullName: repo, CurrentPath: dir, Type: "github_repo", Status: "ready_to_deploy", ImageName: name, Envs: []dto.Env{{Name: "PORT", Value: "80"}}, ExposeSettings: BuildExposeSettings(ds, name, true, "80"), ServiceRuntime: domain.ServiceRuntime{ProjectID: id}}
	// Discover the repository's actual default branch, rather than assuming main.
	if branch, err := ds.GitAdapter.GetBranch(dir); err != nil {
		return app, err
	} else {
		app.Branch = branch
	}
	services := []domain.Service{}
	if ds.FilesystemAdapter.FindDockerComposeFile(dir) {
		configs, err := ds.FilesystemAdapter.GetComposeConfigOfDeploy(dir)
		if err != nil {
			return app, err
		}
		names := []string{}
		for key := range configs {
			names = append(names, key)
		}
		sort.Strings(names)
		buildCount := 0
		for _, key := range names {
			config := configs[key]
			if !ValidID(key) {
				return app, fmt.Errorf("invalid Compose service name %q", key)
			}
			if len(config.DependsOn) > 0 {
				for _, dep := range config.DependsOn {
					if _, ok := configs[dep]; !ok {
						return app, fmt.Errorf("unknown dependency %s", dep)
					}
				}
			}
			if config.HasBuild() {
				buildCount++
				if buildCount > 1 {
					return app, errors.New("Compose import supports one application build; split additional builds into separate repositories")
				}
				dockerfile, err := composeBuild(config.Build)
				if err != nil {
					return app, err
				}
				app.Name = key
				app.Dockerfile = dockerfile
				app.Envs = mergeEnvs(app.Envs, config.GetEnvironmentVariables())
				app.Cmd = config.Cmd
				continue
			}
			if config.Image == "" {
				return app, fmt.Errorf("Compose service %s needs an image", key)
			}
			imageService := domain.Service{Id: utils.GenerateRandomPassword(16), Name: key, Type: "image", ImageName: config.Image, Status: "ready_to_deploy", Envs: config.GetEnvironmentVariables(), Cmd: config.Cmd, ServiceRuntime: domain.ServiceRuntime{ProjectID: id}}
			if len(adapter.DefaultDataTargets(config.Image)) > 0 {
				imageService.Type = "database"
			}
			for _, v := range config.Volumes {
				parts := strings.Split(v, ":")
				if len(parts) != 2 || !ValidID(parts[0]) || !filepath.IsAbs(parts[1]) || strings.Contains(parts[1], "..") {
					return app, fmt.Errorf("Compose service %s supports only named-volume:/absolute/path mounts", key)
				}
				imageService.Volumes = append(imageService.Volumes, domain.ServiceVolume{Name: "justdeploy-" + strings.ToLower(id) + "-" + strings.ToLower(parts[0]), Target: parts[1]})
			}
			// Ports in imported dependency services are not published. Parse only to retain a valid internal port.
			if len(config.Ports) > 0 {
				port := strings.Split(strings.Split(config.Ports[0], "/")[0], ":")
				n, err := strconv.Atoi(port[len(port)-1])
				if err != nil || n < 1 || n > 65535 {
					return app, fmt.Errorf("invalid Compose port for %s", key)
				}
				imageService.ExposeSettings.ExposePort = strconv.Itoa(n)
			}
			if err := ValidateService(imageService); err != nil {
				return app, err
			}
			services = append(services, imageService)
		}
	}
	if err := ValidateService(app); err != nil {
		return app, err
	}
	if err := UpdateServiceCommitInfo(ds, &app); err != nil {
		return app, err
	}
	services = append(services, app)
	if err := ds.DatabaseAdapter.SaveServices(services); err != nil {
		return app, err
	}
	success = true
	return app, nil
}
func composeBuild(value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		if v == "." || v == "./" {
			return "", nil
		}
	case map[string]interface{}:
		for key := range v {
			if key != "context" && key != "dockerfile" {
				return "", fmt.Errorf("unsupported Compose build option %s", key)
			}
		}
		context, _ := v["context"].(string)
		if context != "" && context != "." && context != "./" {
			break
		}
		file, _ := v["dockerfile"].(string)
		if file != "" && (filepath.IsAbs(file) || strings.Contains(file, "..")) {
			break
		}
		return file, nil
	}
	return "", errors.New("Compose application build context must be the repository root")
}
func mergeEnvs(base, extra []dto.Env) []dto.Env {
	for _, env := range extra {
		found := false
		for i := range base {
			if base[i].Name == env.Name {
				base[i] = env
				found = true
			}
		}
		if !found {
			base = append(base, env)
		}
	}
	return base
}
