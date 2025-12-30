package user

import (
	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/domain"
)

func (repo *userdb) UpdateAvatar(id uint, url string) (u *domain.User, err error) {
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

	u = userPo2Domain(&po)
	return
}

func (repo *userdb) UpdateTOTPSecret(uid uint, secret string) error {
	return repo.db.
		Model(&model.UserModel{}).
		Where("id = ?", uid).
		Update("totp_secret", secret).
		Error
}
