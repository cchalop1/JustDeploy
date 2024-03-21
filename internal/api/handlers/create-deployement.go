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
			return c.String(http.StatusBadRequest, "bad request")
		}

		err = application.DeployApplication(deployService, newDeployDto)

		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusOK, dto.ResponseApi{Message: "Application is deploy"})
	}
}
