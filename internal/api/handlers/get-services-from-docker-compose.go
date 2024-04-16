package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func filserviceDto(services []dto.ServiceDto) []string {
	serviceNames := []string{}
	for _, service := range services {
		serviceNames = append(serviceNames, service.Name)
	}
	return serviceNames
}

func GetServicesFromDockerComposeHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		deployId := c.Param("deployId")

		services, err := application.GetServicesFromDockerCompose(deployService, deployId)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, filserviceDto(services))
	}
}
