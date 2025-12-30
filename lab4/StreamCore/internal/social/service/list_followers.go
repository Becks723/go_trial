package service

import (
	"StreamCore/kitex_gen/common"
	"StreamCore/kitex_gen/social"
	"StreamCore/pkg/env"
	"StreamCore/pkg/util"
	"errors"
	"fmt"
)

func (s *SocialService) ListFollowers(query *social.ListFollowersQuery) (*social.SocialList, error) {
	uid, err := util.ParseUint(query.UserId)
	if err != nil {
		return nil, errors.New("bad uid format")
	}

	var limit, page int
	if query.PageSize == nil {
		limit = env.Instance().Social_DefaultPageSize
	} else {
		limit = int(*query.PageSize)
	}

	if query.PageNum == nil {
		page = 0
	} else {
		page = int(*query.PageNum)
	}

	followers, total, err := s.db.QueryFollowers(s.ctx, uid, limit, page)
	if err != nil {
		return nil, fmt.Errorf("err db.QueryFollowers: %w", err)
	}

	data := new(social.SocialList)
	data.Total = int32(total)
	for _, f := range followers {
		data.Items = append(data.Items, s.getSocialInfo(f.TargetUid))
	}
	return data, nil
}

func (s *SocialService) getSocialInfo(uid uint) *common.SocialUserInfo {
	u, err := s.infra.DB.User.GetById(uid)
	if err != nil {
		return &common.SocialUserInfo{
			Id: util.Uint2String(uid),
		}
	} else {
		return &common.SocialUserInfo{
			Id:        util.Uint2String(uid),
			Username:  u.Username,
			AvatarUrl: u.AvatarUrl,
		}
	}

}
