package service

import (
	"StreamCore/biz/domain"
	"StreamCore/biz/model/user"
	"StreamCore/biz/repo"
	"StreamCore/pkg/util"
	"context"
	"errors"
)

type UserService struct {
	repo repo.UserRepo
}

func NewUserService(repo repo.UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (serv *UserService) Register(ctx context.Context, req *user.RegisterReq) (err error) {
	// check potential duplicated username
	if _, err = serv.repo.GetByUsername(req.Username); err == nil {
		err = errors.New("用户名已存在") // TODO: i18n
		return
	}

	// encrypt password
	var encrypted string
	if encrypted, err = util.EncryptPassword(req.Password); err != nil {
		return
	}

	// repo stuff
	do := domain.User{
		Username:  req.Username,
		Password:  encrypted,
		AvatarUrl: "", // TODO: offer default avatar
	}
	if err = serv.repo.Create(&do); err != nil {
		return
	}

	return nil
}
