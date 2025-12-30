package service

import (
	"StreamCore/kitex_gen/common"
	"StreamCore/kitex_gen/video"
	"StreamCore/pkg/util"
	"fmt"
)

func (s *VideoService) Visit(query *video.VisitQuery) (*common.VideoInfo, error) {
	vid, err := util.ParseUint(query.VideoId)
	if err != nil {
		return nil, fmt.Errorf("bad videoId format: %w", err)
	}

	// get video metadata from db
	v, err := s.db.GetById(vid)
	if err != nil {
		return nil, fmt.Errorf("error db GetById: %w", err)
	}

	// cache OnVisited
	if err = s.cache.OnVisited(s.ctx, vid); err != nil {
		return nil, fmt.Errorf("error cache.OnVisited: %w", err)
	}

	data := &common.VideoInfo{
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
	}
	return data, nil
}
