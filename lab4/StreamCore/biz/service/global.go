package service

import (
	"StreamCore/biz/repo"
	"sync"
)

var (
	userSvc   *UserService
	streamSvc *StreamService

	userOnce   sync.Once
	streamOnce sync.Once
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
