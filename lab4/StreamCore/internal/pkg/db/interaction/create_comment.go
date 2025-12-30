package interaction

import (
	"StreamCore/biz/repo/model"
	"StreamCore/internal/pkg/domain"
	"context"
	"fmt"
)

func (repo *iactiondb) CreateComment(ctx context.Context, c *domain.Comment) (err error) {
	po := comDomain2Po(c)
	if po.ParentId != nil { // is sub, ensure videoId
		var parent model.CommentModel
		err = repo.db.First(&parent, *po.ParentId).Error // call First to throw an error if not found
		if err != nil {
			err = fmt.Errorf("Parent comment(id:%d) not found.", *po.ParentId)
			return
		}
		po.VideoId = parent.VideoId
	}
	return repo.db.
		Model(&model.CommentModel{}).
		Create(&po).
		Error
}
