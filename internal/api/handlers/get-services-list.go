package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"github.com/labstack/echo/v4"
)

func GetServicesList(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		services := make([]dto.ServiceDto, 0)

		// move to config file
		services = append(services, dto.ServiceDto{
			Name:  "Postgresql",
			Image: "https://upload.wikimedia.org/wikipedia/commons/thumb/2/29/Postgresql_elephant.svg/640px-Postgresql_elephant.svg.png",
		})

		services = append(services, dto.ServiceDto{
			Name:  "Redis",
			Image: "https://grafikart.fr/uploads/icons/redis.svg",
		})

		return c.JSON(http.StatusOK, services)
	}
}
