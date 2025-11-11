package repo

import (
	"StreamCore/biz/domain"
	"StreamCore/biz/repo/model"
	"time"

	"gorm.io/gorm"
)

type UserRepo interface {
	Create(u *domain.User) error
	GetByUsername(username string) (u *domain.User, err error)
	GetById(id uint) (u *domain.User, err error)
	UpdateAvatar(id uint, url string) (u *domain.User, err error)
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
	u = po2Domain(&po)
	return u, nil
}

func (repo *UserRepository) GetById(id uint) (u *domain.User, err error) {
	po := model.UserModel{}
	err = repo.db.
		Where("id = ?", id).
		First(&po).
		Error
	if err != nil {
		return nil, err
	}
	u = po2Domain(&po)
	return u, nil
}

func (repo *UserRepository) UpdateAvatar(id uint, url string) (u *domain.User, err error) {
	po := model.UserModel{}
	err = repo.db.
		Where("id = ?", id).
		First(&po).
		Error
	if err != nil {
		return
	}

	po.AvatarUrl = url

	err = repo.db.Save(&po).Error
	if err != nil {
		return
	}

	u = po2Domain(&po)
	return
}

func po2Domain(po *model.UserModel) *domain.User {
	return &domain.User{
		Id:        po.ID,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,
		DeletedAt: deletedAtToPtr(po.DeletedAt),
		Username:  po.Username,
		Password:  po.Password,
		AvatarUrl: po.AvatarUrl,
	}
}

func deletedAtToPtr(t gorm.DeletedAt) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}
