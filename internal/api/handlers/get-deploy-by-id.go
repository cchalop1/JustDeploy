package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"github.com/labstack/echo/v4"
)

func GetDeployByIdHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		deployId := c.Param("id")
		deploy, err := deployService.DatabaseAdapter.GetDeployById(deployId)

		if err != nil {
			return c.JSON(http.StatusNotFound, dto.ResponseApi{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, deploy)
	}
}
