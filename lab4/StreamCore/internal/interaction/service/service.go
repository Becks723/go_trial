package service

import (
	"StreamCore/internal/pkg/base"
	cache "StreamCore/internal/pkg/cache/interaction"
	db "StreamCore/internal/pkg/db/interaction"
	"context"
)

type InteractionService struct {
	ctx   context.Context
	db    db.InteractionDatabase
	cache cache.InteractionCache
	infra *base.InfraSet
}

func NewInteractionService(ctx context.Context, infra *base.InfraSet) *InteractionService {
	return &InteractionService{
		ctx:   ctx,
		db:    infra.DB.Interaction,
		cache: infra.Cache.Interaction,
		infra: infra,
	}
}
