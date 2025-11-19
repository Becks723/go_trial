package model

import (
	"context"
	"fmt"
	"time"
)

type LikeModel struct {
	Id         uint `gorm:"primaryKey"`
	Userid     uint
	TargetId   uint // videoId / commentId
	TargetType int  // 1-video, 2-comment
	Status     int  // 1-like, 2-unlike
	Time       time.Time
}

func (l *LikeModel) WbId() string {
	return fmt.Sprintf("like%d_%d", l.Userid, l.TargetId)
}

func (l *LikeModel) ToWbModel(ctx context.Context) interface{} {
	return l
}
