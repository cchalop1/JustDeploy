package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"github.com/labstack/echo/v4"
)

func GetProjectByIdHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		projectId := c.Param("id")
		project, err := deployService.DatabaseAdapter.GetProjectById(projectId)

		if err != nil {
			return c.JSON(http.StatusNotFound, dto.ResponseApi{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, project)
	}
}
