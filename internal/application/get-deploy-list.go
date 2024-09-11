package application

import (
	"cchalop1.com/deploy/internal/api/http/dto"
	"cchalop1.com/deploy/internal/api/service"
)

func GetDeployList(deployService *service.DeployService) []dto.DeployDto {
	deployList := deployService.DatabaseAdapter.GetDeploys()
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
