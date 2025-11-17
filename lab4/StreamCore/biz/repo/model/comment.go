package model

import "gorm.io/gorm"

type CommentModel struct {
	gorm.Model

	AuthorId   uint
	Author     *UserModel `gorm:"foreignKey:AuthorId;references:ID"`
	VideoId    uint
	Video      *VideoModel `gorm:"foreignKey:VideoId;references:ID"`
	Content    string
	ParentId   *uint // parent comment id, null if root
	LikeCount  int
	ChildCount int
}
