package utils

import (
	"math/rand"
)

const (
	lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits           = "0123456789"
	combined         = lowercaseLetters + uppercaseLetters + digits
)

func GenerateRandomPassword(length int) string {
	if length < 1 {
		return ""
	}

	password := make([]byte, length)

	for i := 0; i < length; i++ {
		password[i] = combined[rand.Intn(len(combined))]
	}

	return string(password)
}
