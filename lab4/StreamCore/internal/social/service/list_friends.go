package service

import (
	"StreamCore/kitex_gen/social"
	"StreamCore/pkg/env"
	"fmt"
)

func (s *SocialService) ListFriends(uid uint, query *social.ListFriendsQuery) (*social.SocialList, error) {
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

	mf, total, err := s.db.QueryMutualFollows(s.ctx, uid, limit, page)
	if err != nil {
		return nil, fmt.Errorf("error db.QueryMutualFollows: %w", err)
	}

	data := new(social.SocialList)
	data.Total = int32(total)
	for _, f := range mf {
		data.Items = append(data.Items, s.getSocialInfo(f.TargetUid))
	}
	return data, nil
}
