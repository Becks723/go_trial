package repo

import (
	"context"
	"time"

	"StreamCore/biz/domain"
	"StreamCore/biz/repo/model"
	redisClient "StreamCore/biz/repo/redis"
)

type SocialRepo interface {
	Create(ctx context.Context, follower, followee uint) error
	Delete(ctx context.Context, follower, followee uint) error
	QueryFollows(ctx context.Context, uid uint, limit, page int) ([]*domain.Follow, int, error)
	QueryFollowers(ctx context.Context, uid uint, limit, page int) ([]*domain.Follow, int, error)
	QueryMutualFollows(ctx context.Context, uid uint, limit, page int) ([]*domain.Follow, int, error)
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

func (repo *SocialRepository) QueryFollows(ctx context.Context, uid uint, limit, page int) (follows []*domain.Follow, total int, err error) {
	var records []*model.FollowModel
	var cnt int64
	tx := repo.db.Where("follower_id = ?", uid).
		Count(&cnt)
	if isPageParamsValid(cnt, limit, page) {
		tx = tx.Limit(limit).
			Offset(limit * page)
	}
	tx.Find(&records)

	for _, po := range records {
		follows = append(follows, followDomain(po))
	}
	total = int(cnt)
	return
}

func (repo *SocialRepository) QueryFollowers(ctx context.Context, uid uint, limit, page int) (follows []*domain.Follow, total int, err error) {
	var records []*model.FollowModel
	var cnt int64
	tx := repo.db.Where("followee_id = ?", uid).
		Count(&cnt)
	if isPageParamsValid(cnt, limit, page) {
		tx = tx.Limit(limit).
			Offset(limit * page)
	}
	tx.Find(&records)

	for _, po := range records {
		follows = append(follows, followerDomain(po))
	}
	total = int(cnt)
	return
}

func (repo *SocialRepository) QueryMutualFollows(ctx context.Context, uid uint, limit, page int) (mf []*domain.Follow, total int, err error) {
	var followers []*model.FollowModel
	repo.db.Where("followee_id = ?", uid).
		Find(&followers)
	for _, po := range followers {
		err := repo.db.Where("follower_id = ? AND followee_id = ?", uid, po.FollowerId).
			First(&model.FollowModel{}).
			Error
		if err != nil {
			continue
		}
		mf = append(mf, &domain.Follow{
			TargetUid: po.FollowerId,
		})
	}
	// cursor
	total = len(mf)
	if isPageParamsValid(int64(total), limit, page) {
		mf = mf[limit*page : limit*(page+1)]
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

func followDomain(po *model.FollowModel) *domain.Follow {
	return &domain.Follow{
		TargetUid: po.FolloweeId,
		StartedAt: po.StartedAt,
	}
}

func followerDomain(po *model.FollowModel) *domain.Follow {
	return &domain.Follow{
		TargetUid: po.FollowerId,
		StartedAt: po.StartedAt,
	}
}
