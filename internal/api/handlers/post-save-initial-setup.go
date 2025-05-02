package handlers

import (
	"net/http"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
)

// PostSaveInitialSetupHandler handles saving the initial setup (API key and domain)
func PostSaveInitialSetupHandler(deployService *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Parse the request body
		setupDto := dto.InitialSetupDto{}
		if err := c.Bind(&setupDto); err != nil {
			return c.JSON(http.StatusBadRequest, dto.ResponseApi{
				Message: "Invalid request format",
			})
		}

		// Save the API key and domain
		err := application.SaveInitialSetup(deployService, setupDto)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ResponseApi{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, dto.ResponseApi{
			Message: "Initial setup completed successfully",
		})
	}
}
