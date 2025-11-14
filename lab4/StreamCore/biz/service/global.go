package service

import (
	"StreamCore/biz/repo"
	"sync"
)

var (
	userSvc   *UserService
	streamSvc *StreamService
	lcSvc     *LikeCommentService

	userOnce   sync.Once
	streamOnce sync.Once
	lcOnce     sync.Once
)

// UserService singleton
func UserSvc() *UserService {
	userOnce.Do(func() {
		userSvc = NewUserService(repo.NewUserRepo())
	})
	return userSvc
}

// StreamService singleton
func StreamSvc() *StreamService {
	streamOnce.Do(func() {
		streamSvc = NewStreamService(repo.NewVideoRepo())
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
