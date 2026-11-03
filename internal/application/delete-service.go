package application

import (
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
)

func DeleteService(ds *service.DeployService, id string) error {
	if !ds.Operation.TryLock() {
		return domain.ErrConflict
	}
	defer ds.Operation.Unlock()
	s, err := ds.DatabaseAdapter.GetServiceById(id)
	if err != nil {
		return err
	}
	if err := ds.DockerAdapter.Delete(s.GetDockerName()); err != nil {
		return err
	}
	if err := ds.DatabaseAdapter.DeleteLogFile(id); err != nil {
		return err
	}
	return ds.DatabaseAdapter.DeleteServiceById(id)
}
