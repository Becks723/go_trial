package social

import (
	"StreamCore/internal/pkg/domain"
	"context"

	"gorm.io/gorm"
)

type SocialDatabase interface {
	Create(ctx context.Context, follower, followee uint) error
	Delete(ctx context.Context, follower, followee uint) error
	QueryFollows(ctx context.Context, uid uint, limit, page int) ([]*domain.Follow, int, error)
	QueryFollowers(ctx context.Context, uid uint, limit, page int) ([]*domain.Follow, int, error)
	QueryMutualFollows(ctx context.Context, uid uint, limit, page int) ([]*domain.Follow, int, error)
}

func NewSocialDatabase(gdb *gorm.DB) SocialDatabase {
	return &socialdb{
		db: gdb,
	}
}

type socialdb struct {
	db *gorm.DB
}
