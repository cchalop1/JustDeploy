package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func PostGithubEvent(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		event := dto.GithubEvent{}
		err := c.Bind(&event)
		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		err = application.ManageGithubEvent(deployService, event)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ResponseApi{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, dto.ResponseApi{Message: "Ok"})
	}
}
