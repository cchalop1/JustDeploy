package application

import (
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
)

func storeBuildLog(ds *service.DeployService, id string, logs domain.LogCollector) error {
	return ds.DatabaseAdapter.SaveLogs(id, logs.GetLogs())
}
