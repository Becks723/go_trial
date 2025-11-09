package service

import (
	"StreamCore/biz/domain"
	"StreamCore/biz/model/common"
	"StreamCore/biz/model/user"
	"StreamCore/biz/repo"
	"StreamCore/pkg/env"
	"StreamCore/pkg/util"
	"context"
	"errors"
	"strconv"
	"time"
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

func (serv *UserService) Login(ctx context.Context, req *user.LoginReq) (data *user.UserInfo, auth *common.AuthenticationInfo, err error) {
	// find user in db
	u, err := serv.repo.GetByUsername(req.Username)
	if err != nil {
		return
	}

	// password correct?
	if !util.CheckPassword(req.Password, u.Password) {
		err = errors.New("密码错误") // TODO: i18n
		return
	}

	// generate access, refresh tokens
	ev := env.Instance()
	atoken, err := util.GenerateAccessToken(u, ev.AccessToken_Secret, serv.hoursOf(ev.AccessToken_ExpiryHours))
	if err != nil {
		return
	}
	rtoken, err := util.GenerateRefreshToken(u, ev.RefreshToken_Secret, serv.hoursOf(ev.RefreshToken_ExpiryHours))
	if err != nil {
		return
	}

	data = &user.UserInfo{
		Id:        strconv.FormatUint(uint64(u.Id), 10),
		CreatedAt: u.CreatedAt.String(),
		UpdatedAt: u.UpdatedAt.String(),
		DeletedAt: serv.timePtrToString(u.DeletedAt),
		Username:  u.Username,
		AvatarUrl: u.AvatarUrl,
	}
	auth = &common.AuthenticationInfo{
		AccessToken:  atoken,
		RefreshToken: rtoken,
	}
	err = nil
	return
}

func (serv *UserService) hoursOf(n int) time.Duration {
	return time.Hour * time.Duration(n)
}

func (serv *UserService) timePtrToString(t *time.Time) string {
	if t != nil {
		return t.String()
	}
	return ""
}
