package application

import (
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
)

func GetServiceList(deployService *service.DeployService) []dto.ServiceDto {
	return deployService.FilesystemAdapter.GetServicesListConfig()
}
