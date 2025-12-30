package social

import (
	redisClient "StreamCore/biz/repo/redis"
	"StreamCore/internal/pkg/db/model"
	"context"
	"time"
)

func (repo *socialdb) Create(ctx context.Context, follower, followee uint) (err error) {
	if !repo.exists(follower, followee) { // no exists, then create
		po := model.FollowModel{
			FollowerId: follower,
			FolloweeId: followee,
			StartedAt:  time.Now(),
		}
		err = repo.db.Create(&po).Error
		if err != nil {
			return
		}

		// cache count
		redisClient.Rdb.Incr(ctx, redisClient.FollowsCountKey(follower))
		redisClient.Rdb.Incr(ctx, redisClient.FollowersCountKey(followee))
	}
	return
}
