package application

import (
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
)

func GetServiceBuildLogs(ds *service.DeployService, id string) ([]domain.Logs, error) {
	if _, err := ds.DatabaseAdapter.GetServiceById(id); err != nil {
		return nil, err
	}
	return ds.DatabaseAdapter.GetLogs(id)
}
