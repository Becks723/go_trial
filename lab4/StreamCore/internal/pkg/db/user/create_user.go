package user

import (
	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/domain"
)

func (repo *userdb) Create(u *domain.User) error {
	po := userDomain2Po(u)
	return repo.db.
		Model(&model.UserModel{}).
		Create(&po).
		Error
}
