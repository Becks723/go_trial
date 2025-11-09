package repo

import (
	"StreamCore/biz/domain"
	"StreamCore/biz/repo/model"

	"gorm.io/gorm"
)

type UserRepo interface {
	Create(u *domain.User) error
	GetByUsername(username string) (u *domain.User, err error)
}

type baseRepository struct {
	db *gorm.DB
}

type UserRepository struct {
	baseRepository
}

func NewUserRepo() *UserRepository {
	return &UserRepository{
		baseRepository{db},
	}
}

func (repo *UserRepository) Create(u *domain.User) error {
	po := model.UserModel{
		Username:  u.Username,
		Password:  u.Password,
		AvatarUrl: u.AvatarUrl,
	}
	return repo.db.
		Model(&model.UserModel{}).
		Create(&po).
		Error
}

func (repo *UserRepository) GetByUsername(username string) (u *domain.User, err error) {
	po := model.UserModel{}
	err = repo.db.
		Where("username = ?", username).
		First(&po).
		Error
	if err != nil {
		return nil, err
	}
	u = &domain.User{
		Username:  po.Username,
		Password:  po.Password,
		AvatarUrl: po.AvatarUrl,
	}
	return u, nil
}
