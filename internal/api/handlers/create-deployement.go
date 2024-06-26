package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func CreateDeployementHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		newDeployDto := dto.NewDeployDto{}

		err := c.Bind(&newDeployDto)

		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ResponseApi{Message: err.Error()})
		}

		if err = c.Validate(newDeployDto); err != nil {
			return c.JSON(http.StatusBadRequest, dto.ResponseApi{Message: err.Error()})
		}

		deploy, err := application.DeployApplication(deployService, newDeployDto)

		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusOK, dto.CreateDeployResponse{Id: deploy.Id})
	}
}
