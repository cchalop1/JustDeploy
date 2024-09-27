package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func DeployProjectHandler(deployService service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var deployProjectDto dto.DeployProjectDto

		if err := c.Bind(&deployProjectDto); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		application.DeployProject(deployProjectDto.ProjectId, deployProjectDto.ServerId)

		return c.JSON(http.StatusOK, "ok")
	}
}
