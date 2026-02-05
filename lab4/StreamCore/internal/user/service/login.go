package service

import (
	"StreamCore/internal/pkg/constants"
	"StreamCore/internal/pkg/pack"
	"StreamCore/kitex_gen/common"
	"StreamCore/kitex_gen/user"
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
	atoken, err := jwt.GenerateAccessToken(u.Id, constants.JWT_AccessSecret, constants.JWT_AccessTokenExpiration)
	if err != nil {
		return nil, nil, fmt.Errorf("failed gen accessToken: %w", err)
	}
	rtoken, err := jwt.GenerateRefreshToken(u.Id, constants.JWT_RefreshSecret, constants.JWT_RefreshTokenExpiration)
	if err != nil {
		return nil, nil, fmt.Errorf("failed gen refreshToken: %w", err)
	}

	auth := &common.AuthenticationInfo{
		AccessToken:  atoken,
		RefreshToken: rtoken,
	}
	return pack.UserInfo(u), auth, nil
}
