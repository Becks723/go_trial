package service

import (
	"StreamCore/kitex_gen/common"
	"StreamCore/kitex_gen/interaction"
	"StreamCore/pkg/util"
	"errors"
	"fmt"
)

func (s *InteractionService) ListLikedVideos(query *interaction.ListLikeQuery) (*interaction.ListLikeRespData, error) {
	uid, err := util.ParseUint(query.UserId)
	if err != nil {
		return nil, fmt.Errorf("bad uid format")
	}

	vids, err := s.cache.GetUserLikedVideos(s.ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("error cache.GetUserLikedVideos: %w", err)
	}

	var videos []*common.VideoInfo
	var failCount int
	for _, vid := range vids {
		v, err := s.infra.DB.Video.GetById(vid)
		if err != nil { // if fail to fetch data, fill in with only vid
			failCount++
			videos = append(videos, &common.VideoInfo{
				Id: util.Uint2String(vid),
			})
		} else {
			videos = append(videos, &common.VideoInfo{
				CreatedAt:    v.CreatedAt.String(),
				UpdatedAt:    v.UpdatedAt.String(),
				DeletedAt:    util.TimePtr2String(v.DeletedAt),
				Id:           util.Uint2String(v.Id),
				UserId:       util.Uint2String(v.AuthorId),
				VideoUrl:     v.VideoUrl,
				CoverUrl:     v.CoverUrl,
				Title:        v.Title,
				Description:  v.Description,
				VisitCount:   int32(v.VisitCount),
				LikeCount:    int32(v.LikeCount),
				CommentCount: int32(v.CommentCount),
			})
		}
	}

	if failCount == len(vids) {
		return nil, errors.New("failed to fetch all video data, something might go wrong")
	}
	data := new(interaction.ListLikeRespData)
	data.Items = videos
	return data, nil
}
