package repo

import (
	"StreamCore/biz/domain"
	"StreamCore/biz/repo/model"
	redisClient "StreamCore/biz/repo/redis"
	"StreamCore/biz/repo/wb"
	"StreamCore/pkg/util"
	"context"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LikeCommentRepo interface {
	LikeVideo(ctx context.Context, uid, vid uint, status int) error
	ListVideoLikes(ctx context.Context, uid uint, limit, page int) ([]*domain.Video, error)
	LikeComment(ctx context.Context, uid, cid uint, status int) error
	CreateComment(ctx context.Context, c *domain.Comment) error
	GetCommentById(cid uint) (*domain.Comment, error)
}

type LcRepository struct {
	baseRepository
}

func NewLikeCommentRepo() LikeCommentRepo {
	return NewLcRepository()
}

func NewLcRepository() *LcRepository {
	return &LcRepository{
		baseRepository{db: db},
	}
}

func (repo *LcRepository) LikeVideo(ctx context.Context, uid, vid uint, status int) (err error) {
	// write fast to cache
	if status == 1 {
		err = redisClient.Rdb.SAdd(ctx, redisClient.VideoLikeKey(vid), uid).Err()
		if err != nil {
			return
		}
		err = redisClient.Rdb.SAdd(ctx, redisClient.UserLikeVidKey(uid), vid).Err()
	} else if status == 2 {
		err = redisClient.Rdb.SRem(ctx, redisClient.VideoLikeKey(vid), uid).Err()
		if err != nil {
			return
		}
		err = redisClient.Rdb.SRem(ctx, redisClient.UserLikeVidKey(uid), vid).Err()
	} else {
		err = fmt.Errorf("Unknown status value: %d", status)
	}
	if err != nil {
		return
	}

	// async write to db
	wbc := likeWbc() // write-behind caching
	err = wbc.Enqueue(ctx, &model.LikeModel{
		Userid:     uid,
		TargetId:   vid,
		TargetType: 1,
		Status:     status,
		Time:       time.Now(),
	})
	return
}

func (repo *LcRepository) ListVideoLikes(ctx context.Context, uid uint, limit, page int) (videos []*domain.Video, err error) {
	raw, err := redisClient.Rdb.SMembers(ctx, redisClient.UserLikeVidKey(uid)).Result()
	if err != nil {
		return
	}

	// cursor
	if isPageParamsValid(int64(len(raw)), limit, page) {
		raw = raw[limit*page : limit*(page+1)]
	}

	vidRepo := NewVideoRepo()
	for _, s := range raw {
		vid := util.String2Uint(s)
		var v *domain.Video
		if v, err = vidRepo.GetById(vid); err != nil {
			return
		}
		videos = append(videos, v)
	}
	return
}

func (repo *LcRepository) LikeComment(ctx context.Context, uid, cid uint, status int) (err error) {
	if status == 1 {
		err = redisClient.Rdb.SAdd(ctx, redisClient.CommentLikeKey(cid), uid).Err()
	} else if status == 2 {
		err = redisClient.Rdb.SRem(ctx, redisClient.CommentLikeKey(cid), uid).Err()
	} else {
		err = fmt.Errorf("Unknown status value: %d", status)
	}
	return
}

func (repo *LcRepository) CreateComment(ctx context.Context, c *domain.Comment) (err error) {
	po := comDomain2Po(c)
	return repo.db.
		Model(&model.CommentModel{}).
		Create(&po).
		Error
}

func (repo *LcRepository) GetCommentById(cid uint) (c *domain.Comment, err error) {
	err = repo.db.
		Model(&model.CommentModel{}).
		Where("id = ?", cid).
		First(&c).
		Error
	return
}

func (repo *LcRepository) BatchUpdateLikes(ctx context.Context, batch []*model.LikeModel) error { // batch should not be slice of interface, or gorm won't recognize it
	return repo.db.Model(&model.LikeModel{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&batch).
		Error
}

func comDomain2Po(c *domain.Comment) *model.CommentModel {
	return &model.CommentModel{
		Model: gorm.Model{
			ID:        c.Id,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			DeletedAt: ptrToDeletedAt(c.DeletedAt),
		},
		AuthorId:   c.AuthorId,
		VideoId:    c.VideoId,
		Content:    c.Content,
		ParentId:   c.ParentId,
		LikeCount:  c.LikeCount,
		ChildCount: c.ChildCount,
	}
}

func comPo2Domain(po *model.CommentModel) *domain.Comment {
	return &domain.Comment{
		Id:         po.ID,
		CreatedAt:  po.CreatedAt,
		UpdatedAt:  po.UpdatedAt,
		DeletedAt:  deletedAtToPtr(po.DeletedAt),
		AuthorId:   po.AuthorId,
		VideoId:    po.VideoId,
		Content:    po.Content,
		ParentId:   po.ParentId,
		LikeCount:  po.LikeCount,
		ChildCount: po.ChildCount,
	}
}

var lOnce, cOnce sync.Once
var lwbc, cwbc *wb.Strategy

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
