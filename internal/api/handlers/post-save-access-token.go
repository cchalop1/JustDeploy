package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

func PostSaveAccessTokenHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		installationId := c.Param("installationId")
		if installationId == "" {
			return c.JSON(http.StatusBadRequest, "Installation ID is required")
		}

		err := application.SaveAccessTokenWithInstallationId(deployService, installationId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, "Access token saved successfully")
	}
}
