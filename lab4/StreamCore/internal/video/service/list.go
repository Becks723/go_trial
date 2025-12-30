package service

import (
	"StreamCore/kitex_gen/common"
	"StreamCore/kitex_gen/video"
	"StreamCore/pkg/util"
	"errors"
	"fmt"
)

func (s *VideoService) List(query *video.ListQuery) (*video.ListRespData, error) {
	uid, err := util.ParseUint(query.UserId)
	if err != nil {
		return nil, errors.New("bad uid format")
	}

	limit := int(query.PageSize)
	page := int(query.PageNum)
	videos, total, err := s.db.FetchByUid(uid, limit, page)
	if err != nil {
		return nil, fmt.Errorf("error db.FetchByUid: %w", err)
	}

	data := new(video.ListRespData)
	data.Total = int32(total)
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
