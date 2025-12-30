package service

import (
	"StreamCore/kitex_gen/social"
	"StreamCore/pkg/env"
	"StreamCore/pkg/util"
	"errors"
	"fmt"
)

func (s *SocialService) ListFollows(query *social.ListFollowsQuery) (*social.SocialList, error) {
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

	follows, total, err := s.db.QueryFollows(s.ctx, uid, limit, page)
	if err != nil {
		return nil, fmt.Errorf("err db.QueryFollows: %w", err)
	}

	data := new(social.SocialList)
	data.Total = int32(total)
	for _, f := range follows {
		data.Items = append(data.Items, s.getSocialInfo(f.TargetUid))
	}
	return data, nil
}
