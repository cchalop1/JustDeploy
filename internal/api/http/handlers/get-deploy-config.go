package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/http/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func GetDeployConfigHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		paramsDeployConfig := dto.ParamsDeployConfigDto{}

		err := c.Bind(&paramsDeployConfig)
		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		return c.JSON(http.StatusOK, application.GetDeployConfig(deployService, paramsDeployConfig))
	}
}
