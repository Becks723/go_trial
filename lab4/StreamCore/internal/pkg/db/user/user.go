package user

import (
	"StreamCore/internal/pkg/domain"

	"gorm.io/gorm"
)

type UserDatabase interface {
	Create(u *domain.User) error
	GetByUsername(username string) (u *domain.User, err error)
	GetById(id uint) (u *domain.User, err error)
	UpdateAvatar(id uint, url string) (u *domain.User, err error)
	UpdateTOTPSecret(uid uint, secret string) error
}

func NewUserDataBase(gdb *gorm.DB) UserDatabase {
	return &userdb{
		db: gdb,
	}
}

type userdb struct {
	db *gorm.DB
}
