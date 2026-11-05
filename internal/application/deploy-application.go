package application

import (
	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"errors"
	"fmt"
	"os"
)

func buildServiceDomain(s domain.Service, base string) string {
	if !s.ExposeSettings.IsExposed {
		return ""
	}
	if s.ExposeSettings.SubDomain == "" {
		return base
	}
	return s.ExposeSettings.SubDomain + "." + base
}
func ensureRouter(ds *service.DeployService) error {
	server := ds.DatabaseAdapter.GetServer()
	if !ValidDomain(server.Domain) {
		return errors.New("configure a valid server domain before exposing a service")
	}
	if err := ds.DockerAdapter.PullTreafikImage(); err != nil {
		return err
	}
	return ds.DockerAdapter.RunRouterWithServer(server)
}
func deployOneService(ds *service.DeployService, s domain.Service, base string) error {
	if err := ValidateService(s); err != nil {
		return err
	}
	if s.ExposeSettings.IsExposed {
		if err := ensureRouter(ds); err != nil {
			return err
		}
	}
	switch s.Type {
	case "github_repo":
		info, err := os.Stat(s.CurrentPath)
		if err != nil || !info.IsDir() {
			return errors.New("source workspace is unavailable; redeploy this repository to fetch it again")
		}
		logs := domain.LogCollector{}
		dockerfile, err := adapter.FindDockerfile(s.CurrentPath)
		if err != nil {
			return err
		}
		if s.Dockerfile != "" {
			dockerfile = s.Dockerfile
		}
		if dockerfile != "" {
			s.Dockerfile = dockerfile
			err = ds.DockerAdapter.BuildImage(s, &logs)
		} else {
			err = ds.DockerAdapter.BuildNixpacksImage(s, &logs)
		}
		// Persist only the latest bounded build logs and redact known service credentials.
		for _, env := range s.Envs {
			if len(env.Value) >= 4 {
				logs.Redact(env.Value)
			}
		}
		settings := ds.DatabaseAdapter.GetSettings()
		logs.Redact(settings.GithubToken)
		if logErr := storeBuildLog(ds, s.Id, logs); logErr != nil {
			return errors.Join(err, logErr)
		}
		if err != nil {
			return fmt.Errorf("build failed (previous container retained): %w", err)
		}
	case "database", "llm", "analytics", "image":
		if err := ds.DockerAdapter.EnsureImage(s.ImageName); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported service type %q", s.Type)
	}
	server := ds.DatabaseAdapter.GetServer()
	if err := ds.DockerAdapter.RunImageWithTLS(s, buildServiceDomain(s, base), server.UseHttps); err != nil {
		return err
	}
	if s.Type == "github_repo" {
		if err := UpdateServiceCommitInfo(ds, &s); err != nil {
			return err
		}
	}
	s.Status = "Running"
	s.ExposeSettings.Tls = server.UseHttps
	s.SetUrl(base)
	return ds.DatabaseAdapter.SaveService(s)
}
func DeployApplication(ds *service.DeployService) error {
	return deployServices(ds, "")
}

// DeployServiceGroup deploys the saved workspace and pending dependencies in the
// selected service's connection group without affecting other projects.
func DeployServiceGroup(ds *service.DeployService, id string) error {
	return deployServices(ds, id)
}

func deployServices(ds *service.DeployService, id string) error {
	if !ds.Operation.TryLock() {
		return domain.ErrConflict
	}
	defer ds.Operation.Unlock()
	services, err := ds.DatabaseAdapter.ReadServices()
	if err != nil {
		return err
	}
	group := ""
	if id != "" {
		selected, err := ds.DatabaseAdapter.GetServiceById(id)
		if err != nil {
			return err
		}
		group = selected.ProjectID
		if group == "" {
			group = selected.Id
		}
	}
	server := ds.DatabaseAdapter.GetServer()
	// Start stateful dependencies first. Already running dependencies are never recreated by a global app deploy.
	for _, github := range []bool{false, true} {
		for _, s := range services {
			serviceGroup := s.ProjectID
			if serviceGroup == "" {
				serviceGroup = s.Id
			}
			if id != "" && serviceGroup != group {
				continue
			}
			if (s.Type == "github_repo") != github {
				continue
			}
			running, err := ds.DockerAdapter.IsServiceRunning(s.GetDockerName())
			if err != nil {
				return err
			}
			if running && s.Status != "ready_to_deploy" {
				continue
			}
			if err := deployOneService(ds, s, server.Domain); err != nil {
				return fmt.Errorf("deploy %s: %w", s.Name, err)
			}
		}
	}
	return nil
}
