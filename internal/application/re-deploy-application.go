package application

import (
	"cchalop1.com/deploy/internal"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func ReDeployApplication(ds *service.DeployService, id string) error {
	if !ds.Operation.TryLock() {
		return domain.ErrConflict
	}
	defer ds.Operation.Unlock()
	return redeployService(ds, id)
}
func redeployService(ds *service.DeployService, id string) error {
	s, err := ds.DatabaseAdapter.GetServiceById(id)
	if err != nil {
		return err
	}
	if s.Type != "github_repo" {
		return deployOneService(ds, *s, ds.DatabaseAdapter.GetServer().Domain)
	}
	settings, err := ds.DatabaseAdapter.ReadSettings()
	if err != nil {
		return err
	}
	token := settings.GithubToken
	if token == "" {
		return errors.New("GitHub integration is not configured")
	}
	if !ds.GithubAdapter.IsTokenValid(token) {
		token, err = renewGithubToken(ds, settings.GithubAppSettings)
		if err != nil {
			return err
		}
	}
	dir, err := ds.FilesystemAdapter.CreateWorkspace()
	if err != nil {
		return err
	}
	keep := false
	defer func() {
		if !keep {
			_ = os.RemoveAll(dir)
		}
	}()
	if err := ds.GitAdapter.CloneRepositoryBranch(s.FullName, dir, token, s.Branch); err != nil {
		return err
	}
	previous := s.CurrentPath
	s.CurrentPath = dir
	if err := deployOneService(ds, *s, ds.DatabaseAdapter.GetServer().Domain); err != nil {
		return err
	}
	keep = true
	// Deploy persisted the successful workspace. Previous workspace cleanup is restricted to generated directories.
	workspaces := filepath.Join(internal.JUSTDEPLOY_FOLDER, "workspaces")
	if filepath.Dir(previous) == workspaces && strings.HasPrefix(filepath.Base(previous), "source-") && previous != dir {
		if err := os.RemoveAll(previous); err != nil {
			return err
		}
	}
	return nil
}
