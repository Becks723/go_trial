package service

import (
	"StreamCore/biz/domain"
	"StreamCore/biz/model/common"
	"StreamCore/biz/model/stream"
	"StreamCore/biz/repo"
	"StreamCore/pkg/util"
	"context"
	"mime/multipart"
	"strconv"
	"time"
)

type StreamService struct {
	repo repo.VideoRepo
}

func NewStreamService(repo repo.VideoRepo) *StreamService {
	return &StreamService{
		repo: repo,
	}
}

func (svc *StreamService) GetVideoFeed(ctx context.Context, query *stream.FeedQuery) (data *stream.FeedResp_Data, err error) {
	var after *time.Time
	if query.LatestTime == "" {
		after = nil
	} else {
		var t time.Time
		t, err = parseTIme(query.LatestTime)
		if err != nil {
			return
		}
		after = &t
	}

	videos, err := svc.repo.Fetch(after)
	if err != nil {
		return
	}

	data = new(stream.FeedResp_Data)
	for _, v := range videos {
		data.Items = append(data.Items, streamDomain2Dto(v))
	}
	return
}

func (svc *StreamService) Publish(ctx context.Context, req *stream.PublishReq, fileHeader *multipart.FileHeader) (err error) {
	return
}

func (svc *StreamService) List(ctx context.Context, query *stream.ListQuery) (data *stream.ListResp_Data, err error) {
	return
}

func (svc *StreamService) Popular(ctx context.Context, query *stream.PopularQuery) (data *stream.PopularResp_Data, err error) {
	return
}

func (svc *StreamService) Search(ctx context.Context, query *stream.SearchReq) (data *stream.SearchResp_Data, err error) {
	return
}

func parseTIme(timestamp string) (t time.Time, err error) {
	unix, err := strconv.ParseUint(timestamp, 10, 64)
	if err != nil {
		return
	}
	t = time.UnixMilli(int64(unix))
	return
}

func streamDomain2Dto(v *domain.Video) *common.VideoInfo {
	return &common.VideoInfo{
		CreatedAt:    v.CreatedAt.String(),
		UpdatedAt:    v.UpdatedAt.String(),
		DeletedAt:    util.TimePtr2String(v.DeletedAt),
		Id:           util.Uint2String(v.Id),
		UserId:       util.Uint2String(v.Author.Id),
		VideoUrl:     v.VideoUrl,
		CoverUrl:     v.CoverUrl,
		Title:        v.Title,
		Description:  v.Description,
		VisitCount:   int32(v.VisitCount),
		LikeCount:    int32(v.LikeCount),
		CommentCount: int32(v.CommentCount),
	}
}
