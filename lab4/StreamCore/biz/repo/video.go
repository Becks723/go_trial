package repo

import (
	"StreamCore/biz/domain"
	"StreamCore/biz/repo/model"
	"time"

	"gorm.io/gorm"
)

type VideoRepo interface {
	Create(v *domain.Video) error
	Fetch(after *time.Time) ([]*domain.Video, error)
}

type VideoRepository struct {
	baseRepository
}

func NewVideoRepo() *VideoRepository {
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
