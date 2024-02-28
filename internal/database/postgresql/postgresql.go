package postgresql

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"ms-auth/internal/models"
)

type Storage struct {
	Db *gorm.DB
}

func New(dsn string) (*Storage, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, err
	}
	return &Storage{Db: db}, nil
}
