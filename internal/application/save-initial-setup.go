package application

import (
	"errors"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"golang.org/x/crypto/bcrypt"
)

func SaveInitialSetup(deployService *service.DeployService, setupDto dto.InitialSetupDto) (string, error) {
	settings := deployService.DatabaseAdapter.GetSettings()

	if settings.AdminEmail != "" {
		return "", errors.New("admin account already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(setupDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	jwtSecret, err := GenerateJwtSecret()
	if err != nil {
		return "", err
	}

	settings.AdminEmail = setupDto.Email
	settings.AdminPasswordHash = string(hash)
	settings.JwtSecret = jwtSecret

	if err := deployService.DatabaseAdapter.SaveSettings(settings); err != nil {
		return "", err
	}

	server := deployService.DatabaseAdapter.GetServer()
	server.Domain = setupDto.Domain
	if err := deployService.DatabaseAdapter.SaveServer(server); err != nil {
		return "", err
	}

	return GenerateJWT(settings.AdminEmail, jwtSecret)
}
