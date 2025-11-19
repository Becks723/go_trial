package repo

import (
	"StreamCore/biz/repo/model"
	"StreamCore/biz/repo/wb"
	"context"
	"sync"
	"time"
)

var lOnce sync.Once
var lwbc *wb.Strategy

func likeWbc() *wb.Strategy {
	lOnce.Do(func() {
		lwbc = wb.NewStrategy(&wb.Config{
			Repo:      &rbRepoCoordinator{},
			QueueSize: 50,
			BatchSize: 25,
			Interval:  10 * time.Second,
		})
	})
	return lwbc
}

type rbRepoCoordinator struct {
}

func (c *rbRepoCoordinator) BatchUpdate(ctx context.Context, batch []interface{}) error {
	if len(batch) == 0 {
		return nil
	}
	switch batch[0].(type) {
	case *model.LikeModel:
		likes := make([]*model.LikeModel, len(batch))
		for i, v := range batch {
			likes[i] = v.(*model.LikeModel)
		}
		return NewLcRepository().BatchUpdateLikes(ctx, likes)

	default:
		return nil
	}
}
