package application

import (
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
)

func GetServerInfo(deployService *service.DeployService) dto.Info {
	apiKey := deployService.DatabaseAdapter.GetSettings().ApiKey
	var serverDto dto.ServerDto

	server := deployService.DatabaseAdapter.GetServer()
	serverDto = server.ToServerDto()

	info := dto.Info{
		Version:         GetVersion(),
		FirstConnection: apiKey == "",
		Server:          serverDto,
	}

	return info
}
