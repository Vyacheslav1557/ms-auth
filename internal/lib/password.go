package lib

import (
	"golang.org/x/crypto/bcrypt"
)

func GenerateFromPassword(password string) (string, error) {
	hp, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hp), nil
}

func CompareHashAndPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
