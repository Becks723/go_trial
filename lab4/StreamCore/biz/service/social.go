package service

import (
	"StreamCore/biz/model/common"
	"StreamCore/biz/model/social"
	"StreamCore/biz/repo"
	"StreamCore/pkg/util"
	"context"
	"fmt"
)

type SocialService struct {
	repo repo.SocialRepo
}

func NewSocialService(repo repo.SocialRepo) *SocialService {
	return &SocialService{
		repo: repo,
	}
}

func (svc *SocialService) Follow(ctx context.Context, req *social.FollowReq) (err error) {
	curUid, err := util.RetrieveUserId(ctx)
	if err != nil {
		return
	}
	followee := util.String2Uint(req.ToUserId)

	switch req.ActionType {
	case 0:
		err = svc.repo.Create(ctx, curUid, followee)
	case 1:
		err = svc.repo.Delete(ctx, curUid, followee)
	default:
		err = fmt.Errorf("Unknown follow action type: %d", req.ActionType)
	}
	return
}

func (svc *SocialService) ListFollows(ctx context.Context, query *social.ListFollowsQuery) (data *social.SocialList, err error) {
	uid := util.String2Uint(query.UserId)

	follows, total, err := svc.repo.QueryFollows(ctx, uid, int(query.PageSize), int(query.PageNum))
	if err != nil {
		return
	}

	data = new(social.SocialList)
	data.Total = int32(total)
	for _, f := range follows {
		data.Items = append(data.Items, getSocialInfo(f.TargetUid))
	}
	return
}

func (svc *SocialService) ListFollowers(ctx context.Context, query *social.ListFollowersQuery) (data *social.SocialList, err error) {
	return
}

func (svc *SocialService) ListFriends(ctx context.Context, query *social.ListFriendsQuery) (data *social.SocialList, err error) {
	return
}

func getSocialInfo(uid uint) *common.SocialUserInfo {
	u, err := repo.NewUserRepo().GetById(uid)
	if err != nil {
		return nil
	}
	return &common.SocialUserInfo{
		Id:        util.Uint2String(uid),
		Username:  u.Username,
		AvatarUrl: u.AvatarUrl,
	}
}
