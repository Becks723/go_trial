package interaction

import (
	"StreamCore/internal/pkg/constants"
	"StreamCore/pkg/util"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type InteractionCache interface {
	OnLiked(ctx context.Context, tarType int, uid, tarId uint) error
	OnUnliked(ctx context.Context, tarType int, uid, tarId uint) error
	GetUserLikedVideos(ctx context.Context, uid uint) ([]uint, error)
	SetUserLikedVideos(ctx context.Context, uid uint, vids []uint) error
}

func NewInteractionCache(rdb *redis.Client) InteractionCache {
	return &iacache{
		rdb: rdb,
	}
}

func (c *iacache) OnLiked(ctx context.Context, tarType int, uid, tarId uint) error {
	// 1. incr like count
	// 2. add biz_id to user_id:biz_type set
	// 3. add user_id to biz_id:biz_type set
	// TODO: atomic? lua?

	err := c.rdb.Incr(ctx, c.likeCountKey(tarType, tarId)).Err()
	if err != nil {
		return err
	}
	err = c.rdb.SAdd(ctx, c.userLikesKey(tarType, uid), tarId).Err()
	if err != nil {
		return err
	}
	err = c.rdb.SAdd(ctx, c.likedUsersKey(tarType, tarId), uid).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *iacache) OnUnliked(ctx context.Context, tarType int, uid, tarId uint) error {
	// 1. incr unlike count
	// 2. remove biz_id from user_id:biz_type set
	// 3. remove user_id from biz_id:biz_type set
	// TODO: atomic? lua?

	err := c.rdb.Incr(ctx, c.unlikeCountKey(tarType, tarId)).Err()
	if err != nil {
		return err
	}
	err = c.rdb.SRem(ctx, c.userLikesKey(tarType, uid), tarId).Err()
	if err != nil {
		return err
	}
	err = c.rdb.SRem(ctx, c.likedUsersKey(tarType, tarId), uid).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *iacache) GetUserLikedVideos(ctx context.Context, uid uint) ([]uint, error) {
	members, err := c.rdb.SMembers(ctx, c.userLikesKey(constants.LikeTarType_Video, uid)).Result()
	if err != nil {
		return nil, err
	}

	vids := make([]uint, len(members))
	for i, m := range members {
		vids[i] = util.String2Uint(m)
	}
	return vids, nil
}

func (c *iacache) SetUserLikedVideos(ctx context.Context, uid uint, vids []uint) error {
	key := c.userLikesKey(constants.LikeTarType_Video, uid)
	err := c.rdb.Del(ctx, key).Err()
	if err != nil {
		return nil
	}

	members := make([]interface{}, len(vids))
	for i, vid := range vids {
		members[i] = vid
	}
	err = c.rdb.SAdd(ctx, key, members).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *iacache) likeCountKey(tarType int, tarId uint) string {
	return fmt.Sprintf("like_count:%d:%d", tarId, tarType)
}

func (c *iacache) unlikeCountKey(tarType int, tarId uint) string {
	return fmt.Sprintf("unlike_count:%d:%d", tarId, tarType)
}

func (c *iacache) userLikesKey(tarType int, uid uint) string {
	return fmt.Sprintf("user_likes:%d:%d", uid, tarType)
}

func (c *iacache) likedUsersKey(tarType int, tarId uint) string {
	return fmt.Sprintf("liked_users:%d:%d", tarId, tarType)
}

type iacache struct {
	rdb *redis.Client
}
