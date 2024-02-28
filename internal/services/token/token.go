package tokenservice

import (
	"github.com/golang-jwt/jwt"
	"ms-auth/internal/config"
	"time"
)

const (
	RefreshTokenMaxAge = 30 * time.Minute
	AccessTokenMaxAge  = 5 * time.Minute
)

type RefreshTokenClaims struct {
	jwt.StandardClaims
}

type AccessTokenClaims struct {
	jwt.StandardClaims
}

func NewAccessToken(claims AccessTokenClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString([]byte(config.Cfg().JWT.Secret))
}

func NewRefreshToken(claims RefreshTokenClaims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return refreshToken.SignedString([]byte(config.Cfg().JWT.Secret))
}

func ParseAccessToken(accessToken string) (*AccessTokenClaims, error) {
	parsedAccessToken, err := jwt.ParseWithClaims(accessToken, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg().JWT.Secret), nil
	})

	return parsedAccessToken.Claims.(*AccessTokenClaims), err
}

func ParseRefreshToken(refreshToken string) (*RefreshTokenClaims, error) {
	parsedRefreshToken, err := jwt.ParseWithClaims(refreshToken, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg().JWT.Secret), nil
	})

	return parsedRefreshToken.Claims.(*RefreshTokenClaims), err
}
