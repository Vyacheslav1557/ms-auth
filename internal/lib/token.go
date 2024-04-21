package lib

import (
	"github.com/golang-jwt/jwt"
	"ms-auth/internal/config"
)

func NewToken(claims jwt.StandardClaims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return refreshToken.SignedString([]byte(config.Cfg().JWTSecret))
}

func ParseToken(token string) (*jwt.StandardClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg().JWTSecret), nil
	})
	return parsedToken.Claims.(*jwt.StandardClaims), err
}
