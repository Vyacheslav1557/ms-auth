package auth

import (
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"log/slog"
	"ms-auth/internal/database/postgresql"
	"ms-auth/internal/models"
	passwordservice "ms-auth/internal/services/password"
	tokenservice "ms-auth/internal/services/token"
	"net/http"
	"strconv"
	"time"
)

const (
	RefreshTokenCookieName = "refresh_token"
	AccessTokenCookieName  = "access_token"
)

func Login(log *slog.Logger, storage *postgresql.Storage) http.HandlerFunc {
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// TODO: подумать над информативными ошибками

	//type Response struct {
	//	Error *string `json:"error"`
	//}

	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var user models.User

		tx := storage.Db.Begin()
		err = tx.Where(&models.User{Username: req.Username}).Take(&user).Error
		if err != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if passwordservice.CompareHashAndPassword(user.HashedPassword, req.Password) != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		jti := uuid.New().String()
		rt, err := tokenservice.NewRefreshToken(
			tokenservice.RefreshTokenClaims{
				StandardClaims: jwt.StandardClaims{
					IssuedAt:  time.Now().Unix(),
					ExpiresAt: time.Now().Add(tokenservice.RefreshTokenMaxAge).Unix(),
					Id:        jti,
					Subject:   strconv.Itoa(int(user.ID)),
				},
			},
		)
		if err != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		at, err := tokenservice.NewAccessToken(
			tokenservice.AccessTokenClaims{
				StandardClaims: jwt.StandardClaims{
					IssuedAt:  time.Now().Unix(),
					ExpiresAt: time.Now().Add(tokenservice.AccessTokenMaxAge).Unix(),
					Subject:   strconv.Itoa(int(user.ID)),
				},
			},
		)
		if err != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		user.RefreshTokenJTI = &jti
		tx.Save(&user)
		tx.Commit()
		http.SetCookie(w, &http.Cookie{
			Name:     RefreshTokenCookieName,
			Value:    rt,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})
		http.SetCookie(w, &http.Cookie{
			Name:     AccessTokenCookieName,
			Value:    at,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})
	}
}

func Logout(log *slog.Logger, storage *postgresql.Storage) http.HandlerFunc {
	// TODO: подумать над информативными ошибками

	//type Response struct {
	//	Error *string `json:"error"`
	//}

	return func(w http.ResponseWriter, r *http.Request) {
		rtck, err := r.Cookie(RefreshTokenCookieName)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		rt, err := tokenservice.ParseRefreshToken(rtck.Value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if time.Unix(rt.ExpiresAt, 0).Before(time.Now()) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var user models.User
		tx := storage.Db.Begin()
		err = tx.Where(&models.User{RefreshTokenJTI: &rt.StandardClaims.Id}).Take(&user).Error
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		user.RefreshTokenJTI = nil
		tx.Save(&user)
		tx.Commit()
		http.SetCookie(w, &http.Cookie{
			Name:   RefreshTokenCookieName,
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		http.SetCookie(w, &http.Cookie{
			Name:   AccessTokenCookieName,
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
	}
}

func Refresh(log *slog.Logger, storage *postgresql.Storage) http.HandlerFunc {
	// TODO: подумать над информативными ошибками

	//type Response struct {
	//	Error *string `json:"error"`
	//}

	return func(w http.ResponseWriter, r *http.Request) {
		rtck, err := r.Cookie(RefreshTokenCookieName)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		rt, err := tokenservice.ParseRefreshToken(rtck.Value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var user models.User

		tx := storage.Db.Begin()
		err = tx.Where(&models.User{RefreshTokenJTI: &rt.StandardClaims.Id}).Take(&user).Error
		if err != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		jti := uuid.New().String()
		newrt, err := tokenservice.NewRefreshToken(
			tokenservice.RefreshTokenClaims{
				StandardClaims: jwt.StandardClaims{
					IssuedAt:  time.Now().Unix(),
					ExpiresAt: time.Now().Add(tokenservice.RefreshTokenMaxAge).Unix(),
					Id:        jti,
					Subject:   strconv.Itoa(int(user.ID)),
				},
			},
		)
		if err != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		at, err := tokenservice.NewAccessToken(
			tokenservice.AccessTokenClaims{
				StandardClaims: jwt.StandardClaims{
					IssuedAt:  time.Now().Unix(),
					ExpiresAt: time.Now().Add(tokenservice.AccessTokenMaxAge).Unix(),
					Subject:   strconv.Itoa(int(user.ID)),
				},
			},
		)
		if err != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		user.RefreshTokenJTI = &jti
		tx.Save(&user)
		tx.Commit()
		http.SetCookie(w, &http.Cookie{
			Name:     RefreshTokenCookieName,
			Value:    newrt,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})
		http.SetCookie(w, &http.Cookie{
			Name:     AccessTokenCookieName,
			Value:    at,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})
	}
}
