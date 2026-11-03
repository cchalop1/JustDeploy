package application

import (
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/domain"
	"fmt"
	"net/mail"
)

func UpdateServerTlsSettings(ds *service.DeployService, settings dto.ServerTlsSettings) error {
	if !ds.Operation.TryLock() {
		return domain.ErrConflict
	}
	defer ds.Operation.Unlock()
	if settings.UseHttps {
		email, err := mail.ParseAddress(settings.Email)
		if err != nil || email.Address != settings.Email {
			return fmt.Errorf("valid ACME email required")
		}
	}
	server := ds.DatabaseAdapter.GetServer()
	server.UseHttps = settings.UseHttps
	server.Email = settings.Email
	// Persist desired settings. The router is reconfigured once on the next deployment if its fingerprint changed.
	return ds.DatabaseAdapter.SaveServer(server)
}
