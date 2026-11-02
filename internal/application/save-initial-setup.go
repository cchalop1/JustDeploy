package application

import (
	"cchalop1.com/deploy/internal/api/dto"
	"cchalop1.com/deploy/internal/api/service"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"strings"
)

func SaveInitialSetup(ds *service.DeployService, setup dto.InitialSetupDto) (string, error) {
	email := strings.TrimSpace(setup.Email)
	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email {
		return "", errors.New("a valid email is required")
	}
	if len(setup.Password) < 12 || len(setup.Password) > 72 {
		return "", errors.New("password must contain 12 to 72 bytes")
	}
	if setup.Domain != "" && !ValidDomain(setup.Domain) {
		return "", errors.New("invalid domain")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(setup.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	secret, err := GenerateJwtSecret()
	if err != nil {
		return "", err
	}
	if err := ds.DatabaseAdapter.InitializeAdmin(email, string(hash), secret, setup.Domain); err != nil {
		return "", err
	}
	return GenerateJWT(email, secret)
}
