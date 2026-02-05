package interaction

import (
	"StreamCore/internal/pkg/constants"
	"StreamCore/internal/pkg/db/model"
	"context"
)

func (repo *iactiondb) FetchUserLikedVideos(ctx context.Context, uid uint, limit, page int) ([]uint, error) {
	var likedVids []uint
	repo.db.Model(&model.LikeRelationModel{}).
		Select("target_id").
		Where("uid = ? AND target_type = ? AND status = ?", uid, constants.LikeTarType_Video, constants.LikeAction_Like).
		Limit(limit).
		Offset(limit * page).
		Scan(&likedVids)
	return likedVids, nil
}
