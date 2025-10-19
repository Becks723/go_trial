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

func (repo *UserRepository) FindUserById(uid uint) (*model.UserModel, error) {
	var match model.UserModel
	err := repo.db.
		Model(&model.UserModel{}).
		Where("id = ?", uid).
		First(&match).
		Error
	return &match, err
}

func (repo *UserRepository) FindUserByName(name string) (*model.UserModel, error) {
	var match model.UserModel
	err := repo.db.
		Model(&model.UserModel{}).
		Where("username = ?", name).
		First(&match).
		Error
	return &match, err
}
