package model

import "gorm.io/gorm"

type UserModel struct {
	gorm.Model
	Username string `gorm:"column:username;unique"`
	Password string `gorm:"column:password"`
}
