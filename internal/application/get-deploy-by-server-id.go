package application

import (
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
)

func GetDeployByServerId(deployService *service.DeployService, serverId string) []dto.DeployDto {
	deployList := deployService.DatabaseAdapter.GetDeployByServerId(serverId)
	deployListDto := make([]dto.DeployDto, len(deployList))

	for i, v := range deployList {
		server, err := deployService.DatabaseAdapter.GetServerById(v.ServerId)
		if err != nil {
			continue
		}
		deployListDto[i] = dto.DeployDto{
			Id:           v.Id,
			Name:         v.Name,
			Server:       server.ToServerDto(),
			EnableTls:    v.EnableTls,
			Email:        v.Email,
			Envs:         v.Envs,
			PathToSource: v.PathToSource,
			Status:       v.Status,
			Url:          v.Url,
		}
	}

	return deployListDto
}
