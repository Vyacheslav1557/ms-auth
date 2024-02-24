package random

import (
	"crypto/rand"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(length int) string {
	randomBytes := make([]byte, length)
	rand.Read(randomBytes)
	for i, b := range randomBytes {
		randomBytes[i] = charset[int(b)%len(charset)]
	}
	return string(randomBytes)
}
