package application

import (
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func ValidateService(s domain.Service) error {
	if !ValidID(s.Id) || !ValidID(s.Name) {
		return fmt.Errorf("invalid service identity")
	}
	if s.ProjectID != "" && !ValidID(s.ProjectID) {
		return fmt.Errorf("invalid project ID")
	}
	switch s.Type {
	case "github_repo", "database", "llm", "analytics", "image":
	default:
		return fmt.Errorf("unsupported service type")
	}
	if s.ExposeSettings.SubDomain != "" && !labelPattern.MatchString(s.ExposeSettings.SubDomain) {
		return fmt.Errorf("invalid subdomain")
	}
	if s.ExposeSettings.IsExposed {
		port, err := strconv.Atoi(s.ExposeSettings.ExposePort)
		if err != nil || port < 1 || port > 65535 {
			return fmt.Errorf("invalid exposed port")
		}
	}
	if s.Dockerfile != "" && (filepath.IsAbs(s.Dockerfile) || strings.Contains(s.Dockerfile, "..")) {
		return fmt.Errorf("Dockerfile must be inside the source workspace")
	}
	seen := map[string]bool{}
	for _, env := range s.Envs {
		if !regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`).MatchString(env.Name) || strings.ContainsRune(env.Value, 0) || seen[env.Name] {
			return fmt.Errorf("invalid or duplicate environment variable")
		}
		seen[env.Name] = true
	}
	return nil
}
func UpdateService(ds *service.DeployService, update domain.Service) (domain.Service, error) {
	if !ds.Operation.TryLock() {
		return domain.Service{}, domain.ErrConflict
	}
	defer ds.Operation.Unlock()
	current, err := ds.DatabaseAdapter.GetServiceById(update.Id)
	if err != nil {
		return domain.Service{}, err
	}
	if current.Revision != update.Revision {
		return *current, fmt.Errorf("settings changed; reload before saving: %w", domain.ErrConflict)
	}
	current.Revision++
	current.Envs = update.Envs
	current.ExposeSettings = update.ExposeSettings
	current.ProjectID = update.ProjectID
	current.Status = "ready_to_deploy"
	if err := ValidateService(*current); err != nil {
		return *current, err
	}
	return *current, ds.DatabaseAdapter.SaveService(*current)
}
