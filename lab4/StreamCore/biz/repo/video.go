package repo

import (
	"StreamCore/biz/domain"
	"StreamCore/biz/repo/model"
	redisClient "StreamCore/biz/repo/redis"
	"StreamCore/pkg/util"
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type VideoRepo interface {
	Create(v *domain.Video) error
	GetById(vid uint) (*domain.Video, error)
	Fetch(after *time.Time) ([]*domain.Video, error)
	FetchByUid(uid uint, limit, page int) ([]*domain.Video, int, error)
	IncrVisit(ctx context.Context, vid uint) error
	FetchByVisits(ctx context.Context, limit, page int, reverse bool) ([]*domain.Video, error)
	Search(keywords string, limit, page int, from, to *time.Time, username string) ([]*domain.Video, int, error)
}

func NewVideoRepo() VideoRepo {
	return NewVideoRepository()
}

type VideoRepository struct {
	baseRepository
}

func NewVideoRepository() *VideoRepository {
	return &VideoRepository{
		baseRepository{db: db},
	}
}

func (repo *VideoRepository) Create(v *domain.Video) (err error) {
	po := vidDomain2Po(v)
	return repo.db.
		Model(&model.VideoModel{}).
		Create(po).
		Error
}

func (repo *VideoRepository) Fetch(after *time.Time) (videos []*domain.Video, err error) {
	var records []*model.VideoModel
	if after == nil {
		err = repo.db.Find(&records).Error
	} else {
		err = repo.db.Where("published_at > ?", after).Find(&records).Error
	}
	if err != nil {
		return
	}

	for _, po := range records {
		videos = append(videos, vidPo2Domain(po))
	}
	return
}

func (repo *VideoRepository) FetchByUid(uid uint, limit, page int) (videos []*domain.Video, total int, err error) {
	var records []*model.VideoModel
	var cnt int64
	tx := repo.db.
		Model(&model.VideoModel{}).
		Where("author_id = ?", uid).
		Count(&cnt)
	if err = tx.Error; err != nil {
		return
	}
	total = int(cnt)

	if isPageParamsValid(cnt, limit, page) {
		tx = tx.Limit(limit).
			Offset(limit * page)
	}
	if err = tx.Find(&records).Error; err != nil {
		return
	}

	for _, po := range records {
		videos = append(videos, vidPo2Domain(po))
	}
	return
}

func (repo *VideoRepository) GetById(vid uint) (v *domain.Video, err error) {
	var po *model.VideoModel
	err = repo.db.
		Model(&model.VideoModel{}).
		Where("id = ?", vid).
		First(&po).
		Error
	if err != nil {
		return
	}
	v = vidPo2Domain(po)
	return
}

func (repo *VideoRepository) IncrVisit(ctx context.Context, vid uint) error {
	// incr score
	member := strconv.FormatUint(uint64(vid), 10)
	_, err := redisClient.Rdb.ZIncrBy(ctx, redisClient.VideoRankKey, 1, member).Result()
	if err != nil {
		return err
	}
	// signal async db update
	visitWbc().SetTask(vid, &visitCache{vid: vid})

	return nil
}

func (repo *VideoRepository) FetchByVisits(ctx context.Context, limit, page int, reverse bool) (videos []*domain.Video, err error) {
	if err = repo.ensureVideoRank(ctx); err != nil {
		return
	}

	res, err := redisClient.Rdb.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:   redisClient.VideoRankKey,
		Start: 0,
		Stop:  -1,
		Rev:   reverse,
	}).Result()
	if err != nil {
		return
	}

	if isPageParamsValid(int64(len(res)), limit, page) {
		res, _ = redisClient.Rdb.ZRangeArgs(ctx, redis.ZRangeArgs{
			Key:   redisClient.VideoRankKey,
			Start: limit * page,
			Stop:  limit*(page+1) - 1,
			Rev:   reverse,
		}).Result()
	}
	for _, s := range res {
		vid := util.String2Uint(s)
		var v *domain.Video
		if v, err = repo.GetById(vid); err != nil {
			return
		}
		videos = append(videos, v)
	}
	return
}

func (repo *VideoRepository) Search(keywords string, limit, page int, from, to *time.Time, username string) (videos []*domain.Video, total int, err error) {
	var records []*model.VideoModel

	tx := repo.db.Table("video_models v")
	if keywords != "" {
		tx = tx.Where("title LIKE ? OR description LIKE ?",
			"%"+keywords+"%", "%"+keywords+"%")
	}
	if from != nil {
		tx = tx.Where("published_at > ?", from)
	}
	if to != nil {
		tx = tx.Where("published_at < ?", to)
	}
	if username != "" {
		tx = tx.Joins("JOIN user_models u ON u.id = v.author_id").
			Where("u.username LIKE ?", "%"+username+"%")
	}
	var cnt int64
	if err = tx.Count(&cnt).Error; err != nil {
		return
	}
	if isPageParamsValid(cnt, limit, page) {
		tx = tx.Limit(limit).
			Offset(limit * page)
	}

	if err = tx.Find(&records).Error; err != nil {
		return
	}

	for _, po := range records {
		videos = append(videos, vidPo2Domain(po))
	}
	total = int(cnt)
	return
}

func (repo *VideoRepository) BatchUpdate(ctx context.Context, batch []interface{}) error {
	if len(batch) == 0 {
		return nil
	}

	for _, raw := range batch {
		po := raw.(*model.VideoModel)
		err := repo.db.Model(&model.VideoModel{}).
			Where("id = ?", po.ID).
			Updates(po).
			Error
		if err != nil {
			return err
		}
	}
	return nil
}

func vidDomain2Po(v *domain.Video) *model.VideoModel {
	return &model.VideoModel{
		Model: gorm.Model{
			ID:        v.Id,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			DeletedAt: ptrToDeletedAt(v.DeletedAt),
		},
		AuthorId:     v.AuthorId,
		VideoUrl:     v.VideoUrl,
		CoverUrl:     v.CoverUrl,
		Title:        v.Title,
		Description:  v.Description,
		VisitCount:   v.VisitCount,
		LikeCount:    v.LikeCount,
		CommentCount: v.CommentCount,
		PublishedAt:  v.PublishedAt,
		EditedAt:     v.EditedAt,
	}
}

func vidPo2Domain(po *model.VideoModel) *domain.Video {
	return &domain.Video{
		Id:           po.ID,
		CreatedAt:    po.CreatedAt,
		UpdatedAt:    po.UpdatedAt,
		DeletedAt:    deletedAtToPtr(po.DeletedAt),
		AuthorId:     po.AuthorId,
		VideoUrl:     po.VideoUrl,
		CoverUrl:     po.CoverUrl,
		Title:        po.Title,
		Description:  po.Description,
		VisitCount:   po.VisitCount,
		LikeCount:    po.LikeCount,
		CommentCount: po.CommentCount,
		PublishedAt:  po.PublishedAt,
		EditedAt:     po.EditedAt,
	}
}

// ensureVideoRank init redis videoRank from mysql
func (repo *VideoRepository) ensureVideoRank(ctx context.Context) (err error) {
	exists, _ := redisClient.Rdb.Exists(ctx, redisClient.VideoRankKey).Result()
	if exists != 0 {
		return
	}

	var records []*model.VideoModel
	err = repo.db.Model(&model.VideoModel{}).Find(&records).Error
	if err != nil {
		return
	}

	var members []redis.Z
	for _, po := range records {
		members = append(members, redis.Z{
			Member: po.ID,
			Score:  float64(po.VisitCount),
		})
	}
	err = redisClient.Rdb.ZAdd(ctx, redisClient.VideoRankKey, members...).Err()
	return
}
