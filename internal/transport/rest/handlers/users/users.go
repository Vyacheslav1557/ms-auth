package users

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"ms-auth/internal/database/postgresql"
	"ms-auth/internal/models"
	passwordservice "ms-auth/internal/services/password"
	"net/http"
)

func NewUser(log *slog.Logger, storage *postgresql.Storage) http.HandlerFunc {
	// TODO: подумать над информативными ошибками

	type Response struct {
		Username string `json:"username"`
		Password string `json:"password"`
		//Error    *string `json:"error"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tx := storage.Db.Begin()
		var lastUser models.User
		err := tx.Last(&lastUser).Error
		username := "user"
		if err != nil {
			username += "1"
		} else {
			username += fmt.Sprintf("%d", lastUser.ID+1)
		}
		pw := passwordservice.GenerateRandomPassword()
		hpw, err := passwordservice.GenerateFromPassword(pw)
		if err != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = tx.Create(&models.User{Username: username, HashedPassword: hpw}).Error
		if err != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tx.Commit()
		err = json.NewEncoder(w).Encode(Response{Username: username, Password: pw})
		if err != nil {
			log.Error(err.Error())
			return
		}
	}
}
