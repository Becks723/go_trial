package service

import (
	"StreamCore/kitex_gen/common"
	"StreamCore/kitex_gen/user"
	"StreamCore/pkg/util"
	"errors"
	"fmt"
)

func (s *UserService) GetInfo(query *user.InfoQuery) (*common.UserInfo, error) {
	// convert string id to uint
	uid, err := util.ParseUint(query.UserId)
	if err != nil {
		return nil, errors.New("bad uid format")
	}

	// find user in db
	u, err := s.db.GetById(uid)
	if err != nil {
		return nil, fmt.Errorf("cannot find user (uid=%d)", uid)
	}

	data := &common.UserInfo{
		Id:        util.Uint2String(u.Id),
		CreatedAt: u.CreatedAt.String(),
		UpdatedAt: u.UpdatedAt.String(),
		DeletedAt: util.TimePtr2String(u.DeletedAt),
		Username:  u.Username,
		AvatarUrl: u.AvatarUrl,
	}
	return data, nil
}
