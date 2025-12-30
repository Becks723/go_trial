package social

import (
	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/domain"
	"context"
)

func (repo *socialdb) QueryFollows(ctx context.Context, uid uint, limit, page int) (follows []*domain.Follow, total int, err error) {
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

func (repo *socialdb) QueryFollowers(ctx context.Context, uid uint, limit, page int) (follows []*domain.Follow, total int, err error) {
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

func (repo *socialdb) QueryMutualFollows(ctx context.Context, uid uint, limit, page int) (mf []*domain.Follow, total int, err error) {
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
