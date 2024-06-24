package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func ConnectNewServer(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		connectNewServerDto := dto.ConnectNewServerDto{}

		err := c.Bind(&connectNewServerDto)
		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		newServer := application.CreateServer(deployService, connectNewServerDto)

		return c.JSON(http.StatusOK, newServer)
	}
}
