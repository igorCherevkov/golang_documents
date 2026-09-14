package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateToken() (string, error) {
	buffer := make([]byte, 32)

	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}

	return hex.EncodeToString(buffer), nil
}
