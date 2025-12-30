package video

import (
	redisClient "StreamCore/biz/repo/redis"
	"StreamCore/pkg/util"
	"time"

	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/domain"

	"github.com/redis/go-redis/v9"
)

func (repo *videodb) Fetch(after *time.Time) (videos []*domain.Video, err error) {
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

func (repo *videodb) FetchByUid(uid uint, limit, page int) (videos []*domain.Video, total int, err error) {
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

func (repo *videodb) FetchByVisits(ctx context.Context, limit, page int, reverse bool) (videos []*domain.Video, err error) {
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
