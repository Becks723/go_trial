package service

import (
	cache "StreamCore/biz/repo/cache/user"
	"StreamCore/internal/pkg/base"
	"StreamCore/internal/pkg/db/user"
	"context"
)

type UserService struct {
	ctx   context.Context
	db    user.UserDatabase
	cache cache.UserCache
	infra *base.InfraSet
}

func NewUserService(ctx context.Context, infra *base.InfraSet) *UserService {
	return &UserService{
		ctx:   ctx,
		db:    infra.DB.User,
		infra: infra,
	}
}
