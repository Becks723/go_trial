package service

import (
	"StreamCore/biz/repo"
	"sync"
)

var (
	userSvc *UserService

	userOnce sync.Once
)

// UserService singleton
func UserSvc() *UserService {
	userOnce.Do(func() {
		userSvc = NewUserService(repo.NewUserRepo())
	})
	return userSvc
}
