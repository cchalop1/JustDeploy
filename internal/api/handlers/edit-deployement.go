package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func EditDeployementHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		containerName := c.Param("name")
		editDeployementDto := dto.EditDeployementDto{}
		err := c.Bind(&editDeployementDto)
		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		application.EditDeploy(deployService, containerName, editDeployementDto)

		return c.JSON(http.StatusOK, dto.ResponseApi{Message: "ok"})
	}
}
