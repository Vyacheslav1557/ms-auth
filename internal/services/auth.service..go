package services

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"log/slog"
	"ms-auth/internal/config"
	"ms-auth/internal/lib"
	"ms-auth/internal/models"
	"ms-auth/internal/storage"
	"strconv"
	"time"
)

type Auth struct {
	log     *slog.Logger
	storage *storage.Storage
}

func New(
	log *slog.Logger,
	storage *storage.Storage,
) *Auth {
	return &Auth{
		log:     log,
		storage: storage,
	}
}

func isPassword(str string) bool {
	return len(str) > 0 && len(str) <= 72
}

func isToken(str string) bool {
	return len(str) > 0
}

func isUsername(str string) bool {
	return len(str) > 0
}

var (
	ErrBadCredentials = errors.New("bad credentials")
)

func (a *Auth) Login(username, pwd string) (string, error) {
	var (
		user  models.User
		err   error
		query string
		tkn   string
	)
	if !isUsername(username) || !isPassword(pwd) {
		return "", ErrBadCredentials
	}
	query = "select * from users where username=$1 limit 1"
	err = a.storage.DB.Get(&user, query, username)
	if err != nil {
		return "", ErrBadCredentials
	}
	err = lib.CompareHashAndPassword(user.HashedPassword, pwd)
	if err != nil {
		return "", ErrBadCredentials
	}
	issuedAt := time.Now().Unix()
	query = "update users set is_logged_in=true, last_login_at=$1 where id=$2"
	_, err = a.storage.DB.Exec(query, issuedAt, user.ID)
	if err != nil {
		return "", err
	}
	tkn, err = lib.NewToken(
		jwt.StandardClaims{
			IssuedAt:  issuedAt,
			ExpiresAt: time.Now().Add(config.Cfg().JWTMaxAge).Unix(),
			Subject:   strconv.FormatInt(user.ID, 10),
		},
	)
	if err != nil {
		return "", err
	}
	return tkn, nil
}

func (a *Auth) Refresh(tkn string) (string, error) {
	var (
		err       error
		parsedTkn *jwt.StandardClaims
		user      models.User
		userID    int64
		query     string
	)
	if !isToken(tkn) {
		return "", ErrBadCredentials
	}
	parsedTkn, err = lib.ParseToken(tkn)
	if err != nil {
		return "", ErrBadCredentials
	}
	if time.Unix(parsedTkn.ExpiresAt, 0).Before(time.Now()) {
		return "", ErrBadCredentials
	}
	userID, err = strconv.ParseInt(parsedTkn.Subject, 10, 64)
	if err != nil {
		return "", ErrBadCredentials
	}
	query = "select * from users where id=$1 limit 1"
	err = a.storage.DB.Get(&user, query, userID)
	if err != nil {
		return "", ErrBadCredentials
	}
	if !user.IsLoggedIn {
		return "", ErrBadCredentials
	}
	if time.Unix(user.LastLoginAt, 0).After(time.Unix(parsedTkn.IssuedAt, 0)) {
		return "", ErrBadCredentials
	}
	tkn, err = lib.NewToken(
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(config.Cfg().JWTMaxAge).Unix(),
			Subject:   strconv.FormatInt(userID, 10),
		},
	)
	if err != nil {
		return "", err
	}
	return tkn, nil
}

func (a *Auth) Logout(tkn string) error {
	var (
		err       error
		parsedTkn *jwt.StandardClaims
		user      models.User
		query     string
		userID    int64
	)
	if !isToken(tkn) {
		return ErrBadCredentials
	}
	parsedTkn, err = lib.ParseToken(tkn)
	if err != nil {
		return ErrBadCredentials
	}
	if time.Unix(parsedTkn.ExpiresAt, 0).Before(time.Now()) {
		return ErrBadCredentials
	}
	userID, err = strconv.ParseInt(parsedTkn.Subject, 10, 64)
	if err != nil {
		return ErrBadCredentials
	}
	query = "select * from users where id=$1 limit 1"
	err = a.storage.DB.Get(&user, query, userID)
	if err != nil {
		return ErrBadCredentials
	}
	if !user.IsLoggedIn {
		return ErrBadCredentials
	}
	if time.Unix(user.LastLoginAt, 0).After(time.Unix(parsedTkn.IssuedAt, 0)) {
		return ErrBadCredentials
	}
	query = "update users set is_logged_in=false where id=$1"
	_, err = a.storage.DB.Exec(query, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (a *Auth) Register(username, pwd string) (string, error) {
	var (
		err       error
		hashedPwd string
		query     string
		userID    int64
		tkn       string
	)
	if !isUsername(username) || !isPassword(pwd) {
		return "", ErrBadCredentials
	}
	hashedPwd, err = lib.GenerateFromPassword(pwd)
	if err != nil {
		return "", err
	}
	issuedAt := time.Now().Unix()
	query = "insert into users (username, hashed_password, is_logged_in, last_login_at) values ($1, $2, $3, $4) returning id"
	err = a.storage.DB.QueryRow(query, username, hashedPwd, true, issuedAt).Scan(&userID)
	if err != nil {
		return "", ErrBadCredentials
	}
	tkn, err = lib.NewToken(
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(config.Cfg().JWTMaxAge).Unix(),
			Subject:   strconv.FormatInt(userID, 10),
		},
	)
	if err != nil {
		return "", err
	}
	return tkn, nil
}
