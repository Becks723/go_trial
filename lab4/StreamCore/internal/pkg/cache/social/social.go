package social

import "github.com/redis/go-redis/v9"

type SocialCache interface {
}

func NewSocialCache(rdb *redis.Client) SocialCache {
	return &socialcache{
		rdb: rdb,
	}
}

type socialcache struct {
	rdb *redis.Client
}
