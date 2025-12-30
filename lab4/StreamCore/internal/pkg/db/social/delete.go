package social

import (
	redisClient "StreamCore/biz/repo/redis"
	"StreamCore/internal/pkg/db/model"
	"context"
)

func (repo *socialdb) Delete(ctx context.Context, follower, followee uint) (err error) {
	if repo.exists(follower, followee) { // exists, then delete
		err = repo.db.
			Where("follower_id = ? AND followee_id = ?", follower, followee).
			Delete(&model.FollowModel{}).
			Error
		if err != nil {
			return
		}

		// cache count
		redisClient.Rdb.Decr(ctx, redisClient.FollowsCountKey(follower))
		redisClient.Rdb.Decr(ctx, redisClient.FollowersCountKey(followee))
	}
	return
}
