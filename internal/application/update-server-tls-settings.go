package application

import (
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
)

func UpdateServerTlsSettings(deployService *service.DeployService, tlsSettings dto.ServerTlsSettings) error {
	server := deployService.DatabaseAdapter.GetServer()
	server.UseHttps = tlsSettings.UseHttps
	server.Email = tlsSettings.Email
	err := deployService.DatabaseAdapter.SaveServer(server)
	return err
}
