package handlers

import (
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"github.com/labstack/echo/v4"
	"net/mail"
)

func ServerSettingsHandler(ds *service.DeployService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var settings struct {
			Domain   string `json:"domain"`
			UseHTTPS bool   `json:"useHttps"`
			Email    string `json:"email"`
		}
		if err := c.Bind(&settings); err != nil {
			return echo.NewHTTPError(400, "invalid settings")
		}
		if !application.ValidDomain(settings.Domain) {
			return echo.NewHTTPError(400, "invalid domain")
		}
		if settings.UseHTTPS {
			address, err := mail.ParseAddress(settings.Email)
			if err != nil || address.Address != settings.Email {
				return echo.NewHTTPError(400, "valid ACME email required")
			}
		}
		if !ds.Operation.TryLock() {
			return echo.NewHTTPError(409, "another operation is in progress")
		}
		defer ds.Operation.Unlock()
		server := ds.DatabaseAdapter.GetServer()
		server.Domain = settings.Domain
		server.UseHttps = settings.UseHTTPS
		server.Email = settings.Email
		if err := ds.DatabaseAdapter.SaveServer(server); err != nil {
			return err
		}
		return c.JSON(200, map[string]string{"message": "Server settings saved"})
	}
}
