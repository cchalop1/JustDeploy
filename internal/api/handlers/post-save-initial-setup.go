package handlers

import (
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"cchalop1.com/deploy/internal/domain"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

func PostSaveInitialSetupHandler(ds *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		settings, err := ds.DatabaseAdapter.ReadSettings()
		if err != nil {
			return echo.NewHTTPError(503, "settings unavailable")
		}
		if settings.AdminEmail != "" {
			return echo.NewHTTPError(409, domain.ErrAlreadyConfigured.Error())
		}
		var setup dto.InitialSetupDto
		if err := c.Bind(&setup); err != nil {
			return echo.NewHTTPError(400, "invalid request")
		}
		token, err := application.SaveInitialSetup(ds, setup)
		if errors.Is(err, domain.ErrAlreadyConfigured) {
			return echo.NewHTTPError(409, err.Error())
		}
		if err != nil {
			return echo.NewHTTPError(400, err.Error())
		}
		return c.JSON(http.StatusOK, dto.AuthResponseDto{Token: token})
	}
}
