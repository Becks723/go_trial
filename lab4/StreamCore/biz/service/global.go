package service

import (
	"StreamCore/biz/repo"
	usercache "StreamCore/biz/repo/cache/user"
	"StreamCore/biz/repo/es"
	"sync"
)

var (
	userSvc   *UserService
	streamSvc *StreamService
	lcSvc     *LikeCommentService
	socialSvc *SocialService

	userOnce   sync.Once
	streamOnce sync.Once
	lcOnce     sync.Once
	socialOnce sync.Once
)

// UserService singleton
func UserSvc() *UserService {
	userOnce.Do(func() {
		userSvc = NewUserService(repo.NewUserRepo(), usercache.NewUserCache())
	})
	return userSvc
}

// StreamService singleton
func StreamSvc() *StreamService {
	streamOnce.Do(func() {
		streamSvc = NewStreamService(repo.NewVideoRepo(), es.NewVideoClient())
	})
	return streamSvc
}

// LikeCommentService singleton
func LcSvc() *LikeCommentService {
	lcOnce.Do(func() {
		lcSvc = NewLikeCommentService(repo.NewLikeCommentRepo())
	})
	return lcSvc
}

// SocialService singleton
func SocialSvc() *SocialService {
	socialOnce.Do(func() {
		socialSvc = NewSocialService(repo.NewSocialRepo())
	})
	return socialSvc
}
