package service

import (
	"StreamCore/internal/pkg/pack"
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

	// cache-aside
	// 1. try fetch from cache
	vids, err := s.cache.GetUserLikedVideos(s.ctx, uid)
	if err == nil && len(vids) > 0 {
		return s.saturateVideos(vids)
	}

	// 2. fetch from db, then set back to cache
	vids, err = s.db.FetchUserLikedVideos(s.ctx, uid, int(query.PageSize), int(query.PageNum))
	if err != nil {
		return nil, fmt.Errorf("error db.FetchUserLikedVideos: %w", err)
	}
	err = s.cache.SetUserLikedVideos(s.ctx, uid, vids)
	if err != nil {
		return nil, fmt.Errorf("error cache.SetUserLikedVideos: %w", err) // TODO: log, not return err?
	}
	return s.saturateVideos(vids)
}

func (s *InteractionService) saturateVideos(vids []uint) (*interaction.ListLikeRespData, error) {
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
			videos = append(videos, pack.VideoInfo(v))
		}
	}

	if failCount == len(vids) {
		return nil, errors.New("failed to fetch all video data, something might go wrong")
	}
	data := new(interaction.ListLikeRespData)
	data.Items = videos
	return data, nil
}
