package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"cchalop1.com/deploy/internal/domain"
	"github.com/labstack/echo/v4"
)

func mapServersToDto(servers []domain.Server) []dto.ServerDto {
	serverDtos := make([]dto.ServerDto, len(servers))
	for i, server := range servers {
		serverDtos[i] = dto.ServerDto{
			Id:     server.Id,
			Name:   server.Name,
			Ip:     server.Ip,
			Domain: server.Domain,
			Status: server.Status,
		}
	}
	return serverDtos
}

func GetServerListHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		serverList := application.GetServerList(deployService)
		serverDtoList := mapServersToDto(serverList)
		return c.JSON(http.StatusOK, serverDtoList)
	}
}
