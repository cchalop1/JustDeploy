package application

import (
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
)

func GetDeployByServerId(ds *service.DeployService, id string) []dto.DeployDto {
	if ds.DatabaseAdapter.GetServer().Id != id {
		return []dto.DeployDto{}
	}
	return GetDeployList(ds)
}
