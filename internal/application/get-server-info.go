package application

import (
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
)

func GetServerInfo(deployService *service.DeployService) dto.Info {
	settings := deployService.DatabaseAdapter.GetSettings()
	server := deployService.DatabaseAdapter.GetServer()

	return dto.Info{
		Version:         GetVersion(),
		FirstConnection: settings.AdminEmail == "",
		Server:          server.ToServerDto(),
	}
}
