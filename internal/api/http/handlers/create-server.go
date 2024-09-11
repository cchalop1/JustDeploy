package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/http/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func ConnectNewServer(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		connectNewServerDto := dto.ConnectNewServerDto{}

		err := c.Bind(&connectNewServerDto)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		if connectNewServerDto.Ip == "" {
			return c.JSON(http.StatusBadRequest, dto.ResponseApi{Message: "ip are required"})
		}

		if connectNewServerDto.SshKey == nil && connectNewServerDto.Password == nil {
			return c.JSON(http.StatusBadRequest, dto.ResponseApi{Message: "sshKey or password are required"})
		}

		newServerId, err := application.CreateServer(deployService, connectNewServerDto)

		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ResponseApi{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, dto.ServerDto{Id: newServerId})
	}
}
