package video

import (
	redisClient "StreamCore/biz/repo/redis"
	"context"
	"strconv"
)

func (repo *videodb) IncrVisit(ctx context.Context, vid uint) error {
	// incr score
	member := strconv.FormatUint(uint64(vid), 10)
	_, err := redisClient.Rdb.ZIncrBy(ctx, redisClient.VideoRankKey, 1, member).Result()
	if err != nil {
		return err
	}
	// signal async db update
	visitWbc().SetTask(vid, &visitCache{vid: vid})

	return nil
}
