package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"github.com/labstack/echo/v4"
)

func GetServerByIdHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		serverId := c.Param("id")
		server, err := deployService.DatabaseAdapter.GetServerById(serverId)

		if err != nil {
			return c.JSON(http.StatusNotFound, dto.ResponseApi{Message: err.Error()})
		}

		err = deployService.DockerAdapter.ConnectClient(server)

		if err != nil {
			return c.JSON(http.StatusNotFound, dto.ResponseApi{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, server)
	}
}
