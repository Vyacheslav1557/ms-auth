package models

type User struct {
	ID             int64
	Username       string
	HashedPassword string `db:"hashed_password"`
	IsLoggedIn     bool   `db:"is_logged_in"`
	LastLoginAt    int64  `db:"last_login_at"`
}
