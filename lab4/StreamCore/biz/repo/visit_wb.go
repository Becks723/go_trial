package repo

import (
	"StreamCore/biz/repo/model"
	redisClient "StreamCore/biz/repo/redis"
	"StreamCore/biz/repo/wb"
	"context"
	"strconv"
	"sync"
	"time"

	"gorm.io/gorm"
)

var vOnce sync.Once
var vwbc *wb.DedupStrategy

func visitWbc() *wb.DedupStrategy {
	vOnce.Do(func() {
		vwbc = wb.NewDedupStrategy(&wb.DedupConfig{
			Config: wb.Config{
				Repo:     NewVideoRepository(),
				Interval: 10 * time.Second,
			},
			BatchLimit: 100,
		})
	})
	return vwbc
}

type visitCache struct {
	vid uint
}

func (vc *visitCache) GetCachedValue(ctx context.Context) interface{} {
	// This implementation takes score (i.e. video visits) from redis zset.
	member := strconv.FormatUint(uint64(vc.vid), 10)
	score, _ := redisClient.Rdb.ZScore(ctx, redisClient.VideoRankKey, member).Result()
	return &model.VideoModel{
		Model: gorm.Model{
			ID: vc.vid,
		},
		VisitCount: int(score),
	}
}
