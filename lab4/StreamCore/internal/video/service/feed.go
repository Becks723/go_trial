package service

import (
	"StreamCore/kitex_gen/common"
	"StreamCore/kitex_gen/video"
	"StreamCore/pkg/util"
	"fmt"
	"time"
)

func (s *VideoService) GetVideoFeed(query *video.FeedQuery) (*video.FeedRespData, error) {
	var err error
	var after *time.Time

	if query.LatestTime == nil {
		after = nil
	} else {
		var t time.Time
		t, err = util.FromTimestamp(*query.LatestTime)
		if err != nil {
			return nil, fmt.Errorf("bad timestamp argument: %s", "LatestTime")
		}
		after = &t
	}

	videos, err := s.db.Fetch(after)
	if err != nil {
		return nil, fmt.Errorf("error fetching db video feed: %w", err)
	}

	data := new(video.FeedRespData)
	for _, v := range videos {
		data.Items = append(data.Items, &common.VideoInfo{
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
	return data, nil
}
