package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func CreateServiceHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		deployId := c.Param("deployId")
		serviceName := c.Param("serviceName")

		err := application.CreateService(deployService, deployId, serviceName)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ResponseApi{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, dto.ResponseApi{Message: "Service Is created"})
	}
}
