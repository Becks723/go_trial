package db

import (
	ia "StreamCore/internal/pkg/db/interaction"
	"StreamCore/internal/pkg/db/social"
	"StreamCore/internal/pkg/db/user"
	"StreamCore/internal/pkg/db/video"
	"gorm.io/gorm"
)

type DatabaseSet struct {
	User        user.UserDatabase
	Video       video.VideoDatabase
	Interaction ia.InteractionDatabase
	Social      social.SocialDatabase
}

func NewDatabaseSet(gdb *gorm.DB) *DatabaseSet {
	return &DatabaseSet{
		User:        user.NewUserDataBase(gdb),
		Video:       video.NewVideoDatabase(gdb),
		Interaction: ia.NewInteractionDatabase(gdb),
		Social:      social.NewSocialDatabase(gdb),
	}
}
