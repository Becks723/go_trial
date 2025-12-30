package interaction

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type InteractionCache interface {
	OnVideoLiked(ctx context.Context, vid uint) error
	OnVideoUnliked(ctx context.Context, vid uint) error
	OnCommentLiked(ctx context.Context, cid uint) error
	OnCommentUnliked(ctx context.Context, cid uint) error
	GetUserLikedVideos(ctx context.Context, uid uint) ([]uint, error)
}

func NewInteractionCache(rdb *redis.Client) InteractionCache {
	return &iacache{
		rdb: rdb,
	}
}

func (ia *iacache) OnVideoLiked(ctx context.Context, vid uint) error {

}

func (ia *iacache) OnVideoUnliked(ctx context.Context, vid uint) error {

}

func (ia *iacache) OnCommentLiked(ctx context.Context, vid uint) error {

}

func (ia *iacache) OnCommentUnliked(ctx context.Context, vid uint) error {

}

func (ia *iacache) GetUserLikedVideos(ctx context.Context, uid uint) ([]uint, error) {

}

type iacache struct {
	rdb *redis.Client
}
