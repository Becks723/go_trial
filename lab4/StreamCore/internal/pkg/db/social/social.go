package social

import (
	"StreamCore/internal/pkg/domain"
	"context"
	"time"

	"gorm.io/gorm"
)

type SocialDatabase interface {
	CreateFollow(ctx context.Context, follower, followee uint, time time.Time) error
	UpdateFollowStatus(ctx context.Context, follower, followee uint, status int, time time.Time) error
	GetFollow(ctx context.Context, follower, followee uint) (*domain.Follow, error)
	FetchFollows(ctx context.Context, uid uint, limit, page int) ([]uint, int, error)
	FetchFollowers(ctx context.Context, uid uint, limit, page int) ([]uint, int, error)
	FetchFriends(ctx context.Context, uid uint, limit, page int) ([]uint, int, error)
}

func NewSocialDatabase(gdb *gorm.DB) SocialDatabase {
	return &socialdb{
		db: gdb,
	}
}

type socialdb struct {
	db *gorm.DB
}
