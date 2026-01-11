package model

import "time"

type FollowModel struct {
	Id         uint `gorm:"primaryKey"`
	StartedAt  time.Time
	FollowerId uint
	Follower   *UserModel `gorm:"foreignKey:FollowerId;references:ID"`
	FolloweeId uint
	Followee   *UserModel `gorm:"foreignKey:FolloweeId;references:ID"`
}
