package interaction

import (
	"StreamCore/internal/pkg/domain"
	"context"

	"gorm.io/gorm"
)

type InteractionDatabase interface {
	LikeVideo(ctx context.Context, uid, vid uint, status int) error
	ListVideoLikes(ctx context.Context, uid uint, limit, page int) ([]*domain.Video, error)
	LikeComment(ctx context.Context, uid, cid uint, status int) error
	CreateComment(ctx context.Context, c *domain.Comment) error
	GetCommentById(cid uint) (*domain.Comment, error)
	ListRootComments(vid uint, limit, page int) ([]*domain.Comment, error)
	ListSubComments(cid uint, limit, page int) ([]*domain.Comment, error)
	DeleteCommentById(cid, authorId uint) error
}

func NewInteractionDatabase(gdb *gorm.DB) InteractionDatabase {
	return &iactiondb{
		db: gdb,
	}
}

type iactiondb struct {
	db *gorm.DB
}
