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
		createServiceDto := dto.CreateServiceDto{}

		err := c.Bind(&createServiceDto)

		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		service, err := application.CreateService(deployService, createServiceDto)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ResponseApi{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, service)
	}
}
