package redis

import (
	"context"
	"log"

	"StreamCore/pkg/env"
	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func Init() {
	e := env.Instance()
	Rdb = redis.NewClient(&redis.Options{
		Addr:     e.REDIS_Addr,
		Password: e.REDIS_Password,
		DB:       e.REDIS_DB,
	})

	if err := Rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Error connecting to redis.")
		return
	}
}
