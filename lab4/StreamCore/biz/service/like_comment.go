package service

import (
	"StreamCore/biz/domain"
	"StreamCore/biz/model/comment"
	"StreamCore/biz/model/common"
	"StreamCore/biz/model/like"
	"StreamCore/biz/repo"
	"StreamCore/pkg/util"
	"context"
	"fmt"
)

type LikeCommentService struct {
	repo repo.LikeCommentRepo
}

func NewLikeCommentService(repo repo.LikeCommentRepo) *LikeCommentService {
	return &LikeCommentService{
		repo: repo,
	}
}

func (svc *LikeCommentService) LikeAction(ctx context.Context, req *like.ActionReq) (err error) {
	curUid, err := util.RetrieveUserId(ctx)
	if err != nil {
		return
	}

	if req.VideoId != "" {
		vid := util.String2Uint(req.VideoId)
		err = svc.repo.LikeVideo(ctx, curUid, vid, int(req.ActionType))
	} else if req.CommentId != "" {
		cid := util.String2Uint(req.CommentId)
		err = svc.repo.LikeComment(ctx, curUid, cid, int(req.ActionType))
	}
	return
}

func (svc *LikeCommentService) LikeList(ctx context.Context, query *like.ListQuery) (data *like.ListResp_Data, err error) {
	uid := util.String2Uint(query.UserId)
	videos, err := svc.repo.ListVideoLikes(ctx, uid, int(query.PageSize), int(query.PageNum))
	if err != nil {
		return
	}

	data = new(like.ListResp_Data)
	for _, v := range videos {
		data.Items = append(data.Items, streamDomain2Dto(v))
	}
	return
}

func (svc *LikeCommentService) CommentPublish(ctx context.Context, req *comment.PublishReq) (err error) {
	curUid, err := util.RetrieveUserId(ctx)
	if err != nil {
		return
	}

	var vid uint
	var parentId *uint
	if req.VideoId != "" {
		vid = util.String2Uint(req.VideoId)
	} else if req.CommentId != "" {
		tmp := util.String2Uint(req.CommentId)
		parentId = &tmp
	} else {
		err = fmt.Errorf("Either video_id or comment_id needs to be specified.")
	}

	// db create
	c := &domain.Comment{
		AuthorId: curUid,
		VideoId:  vid,
		ParentId: parentId,
		Content:  req.Content,
	}
	if err = svc.repo.CreateComment(ctx, c); err != nil {
		return
	}
	return
}

func (svc *LikeCommentService) CommentList(ctx context.Context, query *comment.ListQuery) (data *comment.ListResp_Data, err error) {
	var comments []*domain.Comment
	limit, page := int(query.PageSize), int(query.PageNum)

	if query.VideoId != "" {
		vid := util.String2Uint(query.VideoId)
		comments, err = svc.repo.ListRootComments(vid, limit, page)
	} else if query.CommentId != "" {
		cid := util.String2Uint(query.CommentId)
		comments, err = svc.repo.ListSubComments(cid, limit, page)
	} else {
		err = fmt.Errorf("Either video_id or comment_id needs to be specified.")
	}
	if err != nil {
		return
	}

	data = new(comment.ListResp_Data)
	for _, c := range comments {
		data.Items = append(data.Items, comDomain2Dto(c))
	}
	return
}

func (svc *LikeCommentService) CommentDelete(ctx context.Context, req *comment.DeleteReq) (err error) {
	curUid, err := util.RetrieveUserId(ctx)
	if err != nil {
		return
	}

	cid := util.String2Uint(req.CommentId)
	err = svc.repo.DeleteCommentById(cid, curUid)
	if err != nil {
		return
	}
	return
}

func comDomain2Dto(c *domain.Comment) *common.CommentInfo {
	return &common.CommentInfo{
		CreatedAt:  c.CreatedAt.String(),
		UpdatedAt:  c.UpdatedAt.String(),
		DeletedAt:  util.TimePtr2String(c.DeletedAt),
		Id:         util.Uint2String(c.Id),
		UserId:     util.Uint2String(c.AuthorId),
		VideoId:    util.Uint2String(c.VideoId),
		ParentId:   util.Uint2StringOrEmpty(c.ParentId),
		Content:    c.Content,
		LikeCount:  int32(c.LikeCount),
		ChildCount: int32(c.ChildCount),
	}
}
