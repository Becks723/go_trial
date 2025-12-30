package interaction

import (
	redisClient "StreamCore/biz/repo/redis"
	"StreamCore/internal/pkg/domain"
	"StreamCore/pkg/util"
	"context"
)

func (repo *iactiondb) ListVideoLikes(ctx context.Context, uid uint, limit, page int) (videos []*domain.Video, err error) {
	raw, err := redisClient.Rdb.SMembers(ctx, redisClient.UserLikeVidKey(uid)).Result()
	if err != nil {
		return
	}

	// cursor
	if isPageParamsValid(int64(len(raw)), limit, page) {
		raw = raw[limit*page : limit*(page+1)]
	}

	vidRepo := NewVideoRepo()
	for _, s := range raw {
		vid := util.String2Uint(s)
		var v *domain.Video
		if v, err = vidRepo.GetById(vid); err != nil {
			return
		}
		videos = append(videos, v)
	}
	return
}
