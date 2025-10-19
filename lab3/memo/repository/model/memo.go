package model

import (
	"time"
)

type MemoModel struct {
	Id        uint `gorm:"primaryKey"`
	Title     string
	Content   string
	Status    MemoStatus `gorm:"default 0"`
	CreatedAt *time.Time
	StartsAt  *time.Time
	EndsAt    *time.Time
	Uid       uint
	User      UserModel `gorm:"foreignKey:Uid;references:ID"`
}

type MemoStatus int

const (
	MemoStatusPending    MemoStatus = iota // 未开始
	MemoStatusProcessing                   // 进行中
	MemoStatusCompleted                    // 已完成
)
