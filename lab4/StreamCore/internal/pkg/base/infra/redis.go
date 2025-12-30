package infra

import (
	"StreamCore/pkg/env"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func InitRedis() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     env.Instance().REDIS_Addr,
		Password: env.Instance().REDIS_Password,
		DB:       env.Instance().REDIS_DB,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("error connecting to redis: %w", err)
	}

	return rdb, nil
}
