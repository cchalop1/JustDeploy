package application

import (
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
)

func GetDeployList(ds *service.DeployService) []dto.DeployDto {
	result := []dto.DeployDto{}
	server := ds.DatabaseAdapter.GetServer()
	for _, s := range ds.DatabaseAdapter.GetServices() {
		result = append(result, dto.DeployDto{Id: s.Id, Name: s.Name, Server: server.ToServerDto(), EnableTls: s.ExposeSettings.Tls, PathToSource: s.CurrentPath, Envs: s.Envs, Status: s.Status, Url: s.Url})
	}
	return result
}
