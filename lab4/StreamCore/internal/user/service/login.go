package service

import (
	"StreamCore/kitex_gen/common"
	"StreamCore/kitex_gen/user"
	"StreamCore/pkg/env"
	"StreamCore/pkg/util"
	"StreamCore/pkg/util/jwt"
	"errors"
	"fmt"
)

func (s *UserService) Login(req *user.LoginReq) (*common.UserInfo, *common.AuthenticationInfo, error) {
	var err error

	// find user in db
	u, err := s.db.GetByUsername(req.Username)
	if err != nil {
		return nil, nil, errors.New("用户不存在")
	}

	// password correct?
	if !util.CheckPassword(req.Password, u.Password) {
		return nil, nil, errors.New("密码错误")
	}

	// generate access, refresh tokens
	ev := env.Instance()
	atoken, err := jwt.GenerateAccessToken(u.Id, ev.AccessToken_Secret, jwt.HoursOf(ev.AccessToken_ExpiryHours))
	if err != nil {
		return nil, nil, fmt.Errorf("failed gen accessToken: %w", err)
	}
	rtoken, err := jwt.GenerateRefreshToken(u.Id, ev.RefreshToken_Secret, jwt.HoursOf(ev.RefreshToken_ExpiryHours))
	if err != nil {
		return nil, nil, fmt.Errorf("failed gen refreshToken: %w", err)
	}

	data := &common.UserInfo{
		Id:        util.Uint2String(u.Id),
		CreatedAt: u.CreatedAt.String(),
		UpdatedAt: u.UpdatedAt.String(),
		DeletedAt: util.TimePtr2String(u.DeletedAt),
		Username:  u.Username,
		AvatarUrl: u.AvatarUrl,
	}
	auth := &common.AuthenticationInfo{
		AccessToken:  atoken,
		RefreshToken: rtoken,
	}
	return data, auth, nil
}
