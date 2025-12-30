package interaction

import (
	"StreamCore/biz/repo/model"
	"StreamCore/internal/pkg/domain"
)

func (repo *iactiondb) GetCommentById(cid uint) (c *domain.Comment, err error) {
	err = repo.db.
		Model(&model.CommentModel{}).
		Where("id = ?", cid).
		First(&c).
		Error
	return
}
