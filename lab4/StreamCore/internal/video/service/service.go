package service

import (
	"StreamCore/internal/pkg/base"
	cache "StreamCore/internal/pkg/cache/video"
	db "StreamCore/internal/pkg/db/video"
	es "StreamCore/internal/pkg/es/video"
	"context"
)

type VideoService struct {
	ctx   context.Context
	db    db.VideoDatabase
	cache cache.VideoCache
	es    es.VideoElastic
	infra *base.InfraSet
}

func NewVideoService(ctx context.Context, infra *base.InfraSet) *VideoService {
	return &VideoService{
		ctx:   ctx,
		db:    infra.DB.Video,
		cache: infra.Cache.Video,
		es:    es.NewVideoElastic(infra.ES),
		infra: infra,
	}
}
