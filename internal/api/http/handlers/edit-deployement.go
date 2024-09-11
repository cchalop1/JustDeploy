package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/http/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func EditDeployementHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		editDeployDto := dto.EditDeployDto{}
		err := c.Bind(&editDeployDto)

		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		application.EditDeploy(deployService, editDeployDto)

		return c.JSON(http.StatusOK, dto.ResponseApi{Message: "ok"})
	}
}
