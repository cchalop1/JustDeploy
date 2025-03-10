package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

// @Summary Save GitHub access token
// @Description Saves a GitHub access token using the provided installation ID
// @Tags github
// @Accept json
// @Produce json
// @Param installationId path string true "GitHub installation ID"
// @Success 200 {string} string "Access token saved successfully"
// @Failure 400 {string} string "Installation ID is required"
// @Failure 500 {string} string "Error message"
// @Router /api/v1/github/save-token/{installationId} [post]
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
