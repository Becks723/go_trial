package service

import (
	"StreamCore/biz/domain"
	"StreamCore/biz/model/common"
	"StreamCore/biz/model/stream"
	"StreamCore/biz/repo"
	"StreamCore/pkg/util"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"mime/multipart"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
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

func (svc *StreamService) Publish(ctx context.Context, req *stream.PublishReq, videoHeader, coverHeader *multipart.FileHeader) (err error) {
	var (
		localPrefix  = "./uploads"
		accessPrefix = "/static"
	)
	curUid, err := util.RetrieveUserId(ctx)
	if err != nil {
		err = errors.New("Error retrieving user info.")
		return
	}

	if !isValidVideo(videoHeader) {
		err = errors.New("Bad stream format.")
		return
	}

	// save video locally
	dir := fmt.Sprintf(localPrefix + accessPrefix + "/videos/")
	name := uuid.New().String()
	vdst := dir + name + ".mp4"
	if err = saveFile(videoHeader, vdst); err != nil {
		return
	}

	// get cover
	var cdst string
	if coverHeader != nil {
		if !isValidImage(coverHeader) {
			err = errors.New("cover: Bad image format.")
			return
		}

		cdst = dir + name + ".jpg"
		if err = saveFile(coverHeader, cdst); err != nil {
			return
		}
	} else {
		cdst, err = randCover(vdst, dir)
		if err != nil {
			err = fmt.Errorf("Error getting cover.")
			return
		}
	}

	// update db
	vurl, _ := strings.CutPrefix(vdst, localPrefix)
	curl, _ := strings.CutPrefix(cdst, localPrefix)
	now := time.Now()
	v := domain.Video{
		AuthorId:    curUid,
		VideoUrl:    vurl,
		CoverUrl:    curl,
		Title:       req.Title,
		Description: req.Description,
		PublishedAt: now,
		EditedAt:    now,
	}
	if err = svc.repo.Create(&v); err != nil {
		return
	}

	return nil
}

func (svc *StreamService) List(ctx context.Context, query *stream.ListQuery) (data *stream.ListResp_Data, err error) {
	uid, err := util.ParseUint(query.UserId)
	if err != nil {
		err = errors.New("Bad uid format.")
		return
	}

	userRepo := repo.NewUserRepo()
	_, err = userRepo.GetById(uid)
	if err != nil {
		return
	}

	limit := int(query.PageSize)
	page := int(query.PageNum)
	videos, total, err := svc.repo.FetchByUid(uid, limit, page)
	if err != nil {
		return
	}

	data = new(stream.ListResp_Data)
	data.Total = int32(total)
	for _, v := range videos {
		data.Items = append(data.Items, streamDomain2Dto(v))
	}
	return
}

func (svc *StreamService) Popular(ctx context.Context, query *stream.PopularQuery) (data *stream.PopularResp_Data, err error) {
	videos, err := svc.repo.FetchByVisits(ctx, int(query.PageSize), int(query.PageNum), false)
	if err != nil {
		return
	}

	data = new(stream.PopularResp_Data)
	for _, v := range videos {
		data.Items = append(data.Items, streamDomain2Dto(v))
	}
	return
}

func (svc *StreamService) Search(ctx context.Context, query *stream.SearchReq) (data *stream.SearchResp_Data, err error) {
	// resolve from/toDate
	var from, to *time.Time
	var tmp time.Time
	if query.FromDate != "" {
		tmp, err = parseTIme(query.FromDate)
		if err == nil {
			from = &tmp
		}
	}
	if query.ToDate != "" {
		tmp, err = parseTIme(query.ToDate)
		if err == nil {
			to = &tmp
		}
	}

	// core search
	videos, total, err := svc.repo.Search(query.Keywords, int(query.PageSize), int(query.PageNum), from, to, query.Username)
	if err != nil {
		return
	}
	data = new(stream.SearchResp_Data)
	data.Total = int32(total)
	for _, v := range videos {
		data.Items = append(data.Items, streamDomain2Dto(v))
	}
	return
}

func (svc *StreamService) Visit(ctx context.Context, query *stream.VisitQuery) (data *common.VideoInfo, err error) {
	vid, err := util.ParseUint(query.VideoId)
	if err != nil {
		return
	}

	// get video metadata from db
	v, err := svc.repo.GetById(vid)
	if err != nil {
		return
	}
	data = streamDomain2Dto(v)

	// db increase visit
	if err = svc.repo.IncrVisit(ctx, vid); err != nil {
		return
	}

	return
}

func streamDomain2Dto(v *domain.Video) *common.VideoInfo {
	return &common.VideoInfo{
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
}

func randCover(videoPath, coverDir string) (coverPath string, err error) {
	rand.New(rand.NewSource(time.Now().Unix()))

	// get video duration
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		videoPath)
	output, err := cmd.Output()
	if err != nil {
		err = fmt.Errorf("Error retrieving video duration: %s", err.Error())
		return
	}

	// random a timepoint between 20% - 80%
	var duration float64
	fmt.Sscanf(string(output), "%f", &duration)
	if duration <= 0 {
		err = errors.New("Error reading video duration.")
		return
	}
	sec := duration * (0.2 + 0.6*rand.Float64())

	// build cover path
	ext := filepath.Ext(videoPath)
	coverPath = coverDir +
		filepath.Base(videoPath[:len(videoPath)-len(ext)]) + ".jpg"

	// extract frame
	cmd = exec.Command("ffmpeg",
		"-ss", fmt.Sprintf("%.2f", sec),
		"-i", videoPath,
		"-frames:v", "1",
		"-q:v", "2",
		coverPath)
	if err = cmd.Run(); err != nil {
		err = fmt.Errorf("Error writing to cover: %s", err.Error())
		return
	}
	return coverPath, nil
}
