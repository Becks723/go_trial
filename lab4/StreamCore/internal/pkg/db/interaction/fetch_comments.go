package interaction

import (
	"StreamCore/biz/repo/model"
	"StreamCore/internal/pkg/domain"
)

func (repo *iactiondb) ListRootComments(vid uint, limit, page int) (comments []*domain.Comment, err error) {
	var records []*model.CommentModel
	var cnt int64
	tx := repo.db.
		Model(&model.CommentModel{}).
		Where("video_id = ? AND parent_id IS NULL", vid).
		Count(&cnt)

	if isPageParamsValid(cnt, limit, page) {
		tx = tx.Limit(limit).
			Offset(limit * page)
	}
	if err = tx.Find(&records).Error; err != nil {
		return
	}
	for _, po := range records {
		comments = append(comments, comPo2Domain(po))
	}
	return
}

func (repo *iactiondb) ListSubComments(cid uint, limit, page int) (comments []*domain.Comment, err error) {
	var records []*model.CommentModel
	var cnt int64
	tx := repo.db.
		Model(&model.CommentModel{}).
		Where("parent_id = ?", cid).
		Count(&cnt)

	if isPageParamsValid(cnt, limit, page) {
		tx = tx.Limit(limit).
			Offset(limit * page)
	}
	if err = tx.Find(&records).Error; err != nil {
		return
	}
	for _, po := range records {
		comments = append(comments, comPo2Domain(po))
	}
	return
}
