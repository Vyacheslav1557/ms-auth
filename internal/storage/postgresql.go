package storage

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"os"
)

type Storage struct {
	DB *sqlx.DB
}

func New(dataSourceName string) (*Storage, error) {
	db, err := sqlx.Connect("pgx", dataSourceName)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	return &Storage{
		DB: db,
	}, nil
}
