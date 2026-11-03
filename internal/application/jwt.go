package application

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJwtSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateJWT(email, secret string) (string, error) {
	if email == "" || len(secret) < 32 {
		return "", errors.New("authentication is not configured")
	}
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(8 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenStr, secret string) (string, error) {
	if len(secret) < 32 {
		return "", errors.New("authentication is not configured")
	}
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}
	email, ok := claims["email"].(string)
	if !ok || email == "" || !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return "", errors.New("invalid email claim or missing expiry")
	}
	return email, nil
}
