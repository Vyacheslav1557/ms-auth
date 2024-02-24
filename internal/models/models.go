package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID             uint   `gorm:"primarykey"`
	Username       string `gorm:"not null"`
	HashedPassword string `gorm:"not null"`
	LastLoginAt    int64
}
