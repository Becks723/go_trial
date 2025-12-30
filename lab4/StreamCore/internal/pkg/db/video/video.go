package video

import (
	"StreamCore/internal/pkg/domain"
	"context"
	"time"

	"gorm.io/gorm"
)

type VideoDatabase interface {
	Create(v *domain.Video) error
	GetById(vid uint) (*domain.Video, error)
	Fetch(after *time.Time) ([]*domain.Video, error)
	FetchByUid(uid uint, limit, page int) ([]*domain.Video, int, error)
	IncrVisit(ctx context.Context, vid uint) error
	FetchByVisits(ctx context.Context, limit, page int, reverse bool) ([]*domain.Video, error)
	Search(keywords string, limit, page int, from, to *time.Time, username *string) ([]*domain.Video, int, error)
}

func NewVideoDatabase(gdb *gorm.DB) VideoDatabase {
	return &videodb{
		db: gdb,
	}
}

type videodb struct {
	db *gorm.DB
}
