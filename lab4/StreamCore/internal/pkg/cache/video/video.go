package video

import (
	redisClient "StreamCore/biz/repo/redis"
	"StreamCore/pkg/util"
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type VideoCache interface {
	OnVisited(ctx context.Context, vid uint) error
	VisitRank(ctx context.Context, limit, page int, desc bool) ([]uint, error)
}

func NewVideoCache(rdb *redis.Client) VideoCache {
	return &videocache{
		rdb: rdb,
	}
}

func (c *videocache) OnVisited(ctx context.Context, vid uint) error {
	// incr score
	member := strconv.FormatUint(uint64(vid), 10)
	_, err := c.rdb.ZIncrBy(ctx, c.zVideoRankKey(), 1, member).Result()
	if err != nil {
		return err
	}

	// TODO
	// signal async db update
	visitWbc().SetTask(vid, &visitCache{vid: vid})

	return nil
}

func (c *videocache) VisitRank(ctx context.Context, limit, page int, desc bool) ([]uint, error) {
	res, err := c.rdb.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:   c.zVideoRankKey(),
		Start: limit * page,
		Stop:  limit*(page+1) - 1,
		Rev:   desc,
	}).Result()
	if err != nil {
		return nil, err
	}

	var videos []uint
	for _, s := range res {
		vid := util.String2Uint(s)
		videos = append(videos, vid)
	}
	return videos, nil
}

func (c *videocache) zVideoRankKey() string {
	return "zVideoRank"
}

type videocache struct {
	rdb *redis.Client
}
