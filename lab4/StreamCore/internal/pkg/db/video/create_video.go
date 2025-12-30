package video

import (
	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/domain"
)

func (repo *videodb) Create(v *domain.Video) (err error) {
	po := vidDomain2Po(v)
	return repo.db.
		Model(&model.VideoModel{}).
		Create(po).
		Error
}
