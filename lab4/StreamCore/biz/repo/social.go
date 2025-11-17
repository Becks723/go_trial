package repo

import (
	"StreamCore/biz/repo/model"
	redisClient "StreamCore/biz/repo/redis"
	"context"
	"time"
)

type SocialRepo interface {
	Create(ctx context.Context, follower, followee uint) error
	Delete(ctx context.Context, follower, followee uint) error
}

func NewSocialRepo() SocialRepo {
	return &SocialRepository{
		baseRepository{db: db},
	}
}

type SocialRepository struct {
	baseRepository
}

func (repo *SocialRepository) Create(ctx context.Context, follower, followee uint) (err error) {
	if !repo.exists(follower, followee) { // no exists, then create
		po := model.FollowModel{
			FollowerId: follower,
			FolloweeId: followee,
			StartedAt:  time.Now(),
		}
		err = repo.db.Create(&po).Error
		if err != nil {
			return
		}

		// cache count
		redisClient.Rdb.Incr(ctx, redisClient.FollowsCountKey(follower))
		redisClient.Rdb.Incr(ctx, redisClient.FollowersCountKey(followee))
	}
	return
}

func (repo *SocialRepository) Delete(ctx context.Context, follower, followee uint) (err error) {
	if repo.exists(follower, followee) { // exists, then delete
		err = repo.db.
			Where("follower_id = ? AND followee_id = ?", follower, followee).
			Delete(&model.FollowModel{}).
			Error
		if err != nil {
			return
		}

		// cache count
		redisClient.Rdb.Decr(ctx, redisClient.FollowsCountKey(follower))
		redisClient.Rdb.Decr(ctx, redisClient.FollowersCountKey(followee))
	}
	return
}

func (repo *SocialRepository) exists(follower, followee uint) bool {
	err := repo.db.
		Where("follower_id = ? AND followee_id = ?", follower, followee).
		First(&model.FollowModel{}).
		Error
	return err == nil
}
