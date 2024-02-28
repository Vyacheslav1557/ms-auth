package models

import "gorm.io/gorm"

// FIXME: добавить поле Salt

type User struct {
	gorm.Model
	ID              uint   `gorm:"primarykey"`
	Username        string `gorm:"not null; unique"`
	HashedPassword  string `gorm:"not null"`
	RefreshTokenJTI *string
}
