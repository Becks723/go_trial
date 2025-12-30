package service

import (
	"StreamCore/kitex_gen/interaction"
	"StreamCore/pkg/util"
	"fmt"
)

func (s *InteractionService) PublishLike(uid uint, req *interaction.PublishLikeReq) error {
	if req.VideoId == nil && req.CommentId == nil {
		return fmt.Errorf("either videoId or commentId should have value")
	}
	if req.ActionType != 1 && req.ActionType != 2 {
		return fmt.Errorf("unknown actionType (value=%d)", req.ActionType)
	}

	if req.VideoId != nil {
		vid, err := util.ParseUint(*req.VideoId)
		if err != nil {
			return fmt.Errorf("bad videoId format: %w", err)
		}
		switch req.ActionType {
		case 1:
			if err = s.cache.OnVideoLiked(s.ctx, vid); err != nil {
				return fmt.Errorf("error cache.OnVideoLiked: %w", err)
			}
		case 2:
			if err = s.cache.OnVideoUnliked(s.ctx, vid); err != nil {
				return fmt.Errorf("error cache.OnVideoUnliked: %w", err)
			}
		}

		// TODO: mq update video like
	} else {
		cid, err := util.ParseUint(*req.CommentId)
		if err != nil {
			return fmt.Errorf("bad commentId format: %w", err)
		}
		switch req.ActionType {
		case 1:
			if err = s.cache.OnCommentLiked(s.ctx, cid); err != nil {
				return fmt.Errorf("error cache.OnCommentLiked: %w", err)
			}
		case 2:
			if err = s.cache.OnCommentUnliked(s.ctx, cid); err != nil {
				return fmt.Errorf("error cache.OnCommentUnliked: %w", err)
			}
		}

		// TODO: mq update comment like
	}
	return nil
}
