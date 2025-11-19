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
	s, _ := redisClient.Rdb.Get(ctx, redisClient.VideoVisitCountKey(vc.vid)).Result()
	visits, _ := strconv.ParseInt(s, 10, 32)
	return &model.VideoModel{
		Model: gorm.Model{
			ID: vc.vid,
		},
		VisitCount: int(visits),
	}
}
