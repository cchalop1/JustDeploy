package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"github.com/labstack/echo/v4"
)

func ConnectNewServer(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		connectNewServerDto := dto.ConnectNewServerDto{}

		err := c.Bind(&connectNewServerDto)
		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		return c.JSON(http.StatusOK, true)
	}
}
