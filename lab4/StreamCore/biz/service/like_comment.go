package service

import (
	"StreamCore/biz/domain"
	"StreamCore/biz/model/comment"
	"StreamCore/biz/model/like"
	"StreamCore/biz/repo"
	"StreamCore/pkg/util"
	"context"
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

	var vid, parentId *uint
	var tmp uint
	if req.VideoId != "" {
		tmp = util.String2Uint(req.VideoId)
		vid = &tmp
		// assert video exists
		if _, err = repo.NewVideoRepo().GetById(*vid); err != nil {
			return
		}
	} else if req.CommentId != "" {
		tmp = util.String2Uint(req.CommentId)
		parentId = &tmp
		// assert parent exists
		if _, err = repo.NewLikeCommentRepo().GetCommentById(*parentId); err != nil {
			return
		}
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
	return
}

func (svc *LikeCommentService) CommentDelete(ctx context.Context, req *comment.DeleteReq) (err error) {
	return
}
