package application

import (
	"errors"

	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"golang.org/x/crypto/bcrypt"
)

func Login(deployService *service.DeployService, loginDto dto.LoginDto) (string, error) {
	settings := deployService.DatabaseAdapter.GetSettings()

	if settings.AdminEmail == "" {
		return "", errors.New("no admin account configured")
	}

	if loginDto.Email != settings.AdminEmail {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(settings.AdminPasswordHash), []byte(loginDto.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	return GenerateJWT(settings.AdminEmail, settings.JwtSecret)
}
