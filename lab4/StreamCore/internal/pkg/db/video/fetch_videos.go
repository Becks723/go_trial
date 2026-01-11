package video

import (
	"StreamCore/internal/pkg/db/pack"
	"StreamCore/internal/pkg/db/util"
	"time"

	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/domain"
)

func (repo *videodb) Fetch(after *time.Time) (videos []*domain.Video, err error) {
	var records []*model.VideoModel
	if after == nil {
		err = repo.db.Find(&records).Error
	} else {
		err = repo.db.Where("published_at > ?", after).Find(&records).Error
	}
	if err != nil {
		return nil, err
	}

	for _, po := range records {
		videos = append(videos, repo.packVideo(po))
	}
	return videos, nil
}

func (repo *videodb) FetchByUid(uid uint, limit, page int) (videos []*domain.Video, total int, err error) {
	var records []*model.VideoModel
	var cnt int64
	tx := repo.db.
		Model(&model.VideoModel{}).
		Where("author_id = ?", uid).
		Count(&cnt)
	if err = tx.Error; err != nil {
		return nil, -1, err
	}

	if util.IsPageParamsValid(limit, page) {
		tx = tx.Limit(limit).
			Offset(limit * page)
	}
	if err = tx.Find(&records).Error; err != nil {
		return nil, -1, err
	}

	for _, po := range records {
		videos = append(videos, pack.Video(po))
	}
	return videos, int(cnt), nil
}
