package domain

import "time"

type Video struct {
	Id           uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
	Author       *User
	VideoUrl     string
	CoverUrl     string
	Title        string
	Description  string
	VisitCount   int
	LikeCount    int
	CommentCount int
	PublishedAt  time.Time
	EditedAt     time.Time
}
