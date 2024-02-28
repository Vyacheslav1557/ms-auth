package passwordservice

import (
	"golang.org/x/crypto/bcrypt"
	"ms-auth/internal/lib/random"
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

func GenerateRandomPassword() string {
	return random.RandomString(3) + "-" + random.RandomString(3) + "-" + random.RandomString(3)
}
