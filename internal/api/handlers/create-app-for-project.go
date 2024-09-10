package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func CreateAppForProjectHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		createNewApp := dto.CreateAppDto{}

		err := c.Bind(&createNewApp)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		app, err := application.CreateApp(deployService, createNewApp)

		if err != nil {
			return c.JSON(http.StatusNotFound, dto.ResponseApi{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, app)
	}
}
