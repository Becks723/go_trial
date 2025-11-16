package domain

import "time"

type Comment struct {
	Id         uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
	AuthorId   uint
	VideoId    *uint
	ParentId   *uint
	Content    string
	LikeCount  int
	ChildCount int
}
