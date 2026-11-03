package utils

import (
	"crypto/rand"
	"math/big"
)

const combined = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomPassword(length int) string {
	if length < 1 {
		return ""
	}
	password := make([]byte, length)
	for i := range password {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(combined))))
		if err != nil {
			panic("secure random source unavailable")
		}
		password[i] = combined[n.Int64()]
	}
	return string(password)
}
