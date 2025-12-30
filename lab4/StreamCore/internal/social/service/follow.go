package service

import (
	"StreamCore/kitex_gen/social"
	"StreamCore/pkg/util"
	"errors"
	"fmt"
)

func (s *SocialService) Follow(uid uint, req *social.FollowReq) error {
	followee, err := util.ParseUint(req.ToUserId)
	if err != nil {
		return errors.New("bad uid format")
	}
	if req.ActionType != 0 && req.ActionType != 1 {
		return fmt.Errorf("unknown follow actionType (value=%d)", req.ActionType)
	}

	switch req.ActionType {
	case 0:
		if err = s.db.Create(s.ctx, uid, followee); err != nil {
			return fmt.Errorf("error db.Create: %w", err)
		}
	case 1:
		if err = s.db.Delete(s.ctx, uid, followee); err != nil {
			return fmt.Errorf("error db.Delete: %w", err)
		}
	}
	return nil
}
