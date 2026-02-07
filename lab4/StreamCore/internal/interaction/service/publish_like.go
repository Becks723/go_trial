package service

import (
	"context"
	"fmt"
	"time"

	"StreamCore/internal/pkg/constants"
	"StreamCore/internal/pkg/mq/model"
	"StreamCore/kitex_gen/interaction"
	"StreamCore/pkg/util"
	"github.com/bytedance/sonic"
)

func (s *InteractionService) PublishLike(uid uint, req *interaction.PublishLikeReq) error {
	var err error
	if req.VideoId == nil && req.CommentId == nil {
		return fmt.Errorf("either videoId or commentId should have value")
	}
	if req.ActionType != constants.LikeAction_Like && req.ActionType != constants.LikeAction_Unlike {
		return fmt.Errorf("unknown actionType (value=%d)", req.ActionType)
	}

	var tarType int
	var tarId uint
	if req.VideoId != nil {
		vid, err := util.ParseUint(*req.VideoId)
		if err != nil {
			return fmt.Errorf("bad videoId format: %w", err)
		}
		tarType = constants.LikeTarType_Video
		tarId = vid
	} else {
		cid, err := util.ParseUint(*req.CommentId)
		if err != nil {
			return fmt.Errorf("bad commentId format: %w", err)
		}
		tarType = constants.LikeTarType_Comment
		tarId = cid
	}

	// cache
	switch req.ActionType {
	case 1:
		if err = s.cache.OnLiked(s.ctx, tarType, uid, tarId); err != nil {
			return fmt.Errorf("error cache.OnLiked: %w", err)
		}
	case 2:
		if err = s.cache.OnUnliked(s.ctx, tarType, uid, tarId); err != nil {
			return fmt.Errorf("error cache.OnUnliked: %w", err)
		}
	}

	// mq
	err = s.mq.PublishLikeEvent(s.ctx, &model.LikeEvent{
		TarType: tarType,
		TarId:   tarId,
		Uid:     uid,
		Action:  int(req.ActionType),
		Time:    time.Now(),
	})
	if err != nil {
		return fmt.Errorf("error mq.PublishLikeEvent: %w", err)
	}
	return nil
}

func (s *InteractionService) consumeLike(ctx context.Context) {
	c, err := s.mq.Consumer()
	if err != nil {
		// TODO: log consumer init error
		return
	}
	for {
		msg, err := c.Receive()
		if err != nil {
			break
		}
		var ev model.LikeEvent
		err = sonic.Unmarshal(msg.Body, &ev)
		if err != nil {
			continue
		}

		// >>>> consume like event <<<<
		if ev.Action == constants.LikeAction_Like {
			if err = s.publishLikeToDB(ctx, ev.TarType, ev.Uid, ev.TarId, ev.Time); err != nil { //nolint:staticcheck
				// TODO: log
			}
		} else {
			if err = s.publishUnlikeToDB(ctx, ev.TarType, ev.Uid, ev.TarId, ev.Time); err != nil { //nolint:staticcheck
				// TODO: log
			}
		}
		// =============================

		c.Ack(msg)
	}
}

func (s *InteractionService) publishLikeToDB(ctx context.Context, tarType int, uid, tarId uint, time time.Time) error {
	lastLike, err := s.db.GetLike(ctx, tarType, uid, tarId)
	if err != nil { // not recorded before
		if err = s.db.CreateLike(ctx, tarType, uid, tarId, time); err != nil {
			return fmt.Errorf("error db.CreateLike: %w", err)
		}
		return nil
	}

	if lastLike.Time.Before(time) && // ignore former request
		lastLike.Status == constants.LikeAction_Unlike { // ignore repeated like
		err = s.db.ToggleLikeStatus(ctx, tarType, uid, tarId)
		if err != nil {
			return fmt.Errorf("error db.ToggleLikeStatus: %w", err)
		}
	}
	return nil
}

func (s *InteractionService) publishUnlikeToDB(ctx context.Context, tarType int, uid, tarId uint, time time.Time) error {
	lastLike, err := s.db.GetLike(ctx, tarType, uid, tarId)
	if err != nil { // not recorded before
		return nil
	}

	if lastLike.Time.Before(time) && // ignore former request
		lastLike.Status == constants.LikeAction_Like { // ignore repeated unlike
		err = s.db.ToggleLikeStatus(ctx, tarType, uid, tarId)
		if err != nil {
			return fmt.Errorf("error db.ToggleLikeStatus: %w", err)
		}
	}
	return nil
}
