package repository

import (
	"memo/repository/model"

	"gorm.io/gorm"
)

type baseRepository struct {
	db *gorm.DB
}

type UserRepository struct {
	baseRepository
}

func NewUserRepo() *UserRepository {
	return &UserRepository{
		baseRepository{
			db: db,
		},
	}
}

func (repo *UserRepository) InsertUser(user *model.UserModel) (err error) {
	return repo.db.
		Model(&model.UserModel{}).
		Create(user).
		Error
}

func (repo *UserRepository) FindUserByName(name string) (*model.UserModel, error) {
	var match model.UserModel
	err := repo.db.
		Model(&model.UserModel{}).
		Where("username = ?", name).
		First(&match).
		Error
	if err != nil {
		return nil, err
	}
	return &match, nil
}
