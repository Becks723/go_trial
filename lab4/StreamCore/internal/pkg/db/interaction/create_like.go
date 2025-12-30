package interaction

import (
	"StreamCore/biz/repo/model"
	redisClient "StreamCore/biz/repo/redis"
	"context"
	"fmt"
	"time"
)

func (repo *iactiondb) LikeVideo(ctx context.Context, uid, vid uint, status int) (err error) {
	// write fast to cache
	if status == 1 {
		err = redisClient.Rdb.SAdd(ctx, redisClient.VideoLikeKey(vid), uid).Err()
		if err != nil {
			return
		}
		err = redisClient.Rdb.SAdd(ctx, redisClient.UserLikeVidKey(uid), vid).Err()
	} else if status == 2 {
		err = redisClient.Rdb.SRem(ctx, redisClient.VideoLikeKey(vid), uid).Err()
		if err != nil {
			return
		}
		err = redisClient.Rdb.SRem(ctx, redisClient.UserLikeVidKey(uid), vid).Err()
	} else {
		err = fmt.Errorf("Unknown status value: %d", status)
	}
	if err != nil {
		return
	}

	// async write to db
	wbc := likeWbc() // write-behind caching
	err = wbc.Enqueue(ctx, &model.LikeModel{
		Userid:     uid,
		TargetId:   vid,
		TargetType: 1,
		Status:     status,
		Time:       time.Now(),
	})
	return
}

func (repo *iactiondb) LikeComment(ctx context.Context, uid, cid uint, status int) (err error) {
	if status == 1 {
		err = redisClient.Rdb.SAdd(ctx, redisClient.CommentLikeKey(cid), uid).Err()
	} else if status == 2 {
		err = redisClient.Rdb.SRem(ctx, redisClient.CommentLikeKey(cid), uid).Err()
	} else {
		err = fmt.Errorf("Unknown status value: %d", status)
	}
	return
}
