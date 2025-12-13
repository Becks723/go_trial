package user

import (
	rd "StreamCore/biz/repo/redis"
	"StreamCore/pkg/constants"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func (c *UserCache) SetTOTPPending(ctx context.Context, uid uint, secret string) error {
	cacheKey := c.totpPendingKey(uid)
	return c.rdb.Set(ctx, cacheKey, secret, constants.TOTPSecretExpiry).Err()
}

func (c *UserCache) GetTOTPPending(ctx context.Context, uid uint) (string, error) {
	cacheKey := c.totpPendingKey(uid)
	return c.rdb.Get(ctx, cacheKey).Result()
}

func (c *UserCache) LockCurrentTOTPPeriod(ctx context.Context, uid uint) error {
	timestep := time.Now().UnixMilli() / constants.TOTPInterval
	key := c.totpPeriodKey(uid, timestep)
	return c.rdb.Set(ctx, key, 1, constants.TOTPInterval*2*time.Second).Err()
}

func (c *UserCache) IsCurrentTOTPPeriodLocked(ctx context.Context, uid uint) bool {
	timestep := time.Now().UnixMilli() / constants.TOTPInterval
	key := c.totpPeriodKey(uid, timestep)
	ok, _ := c.rdb.Exists(ctx, key).Result()
	return ok == 1
}

func (c *UserCache) IncreaseTOTPFailure(ctx context.Context, uid uint) error {
	key := c.totpFailureKey(uid)
	ok, err := c.rdb.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if ok == 0 { // key not exists, create key with ttl
		err = c.rdb.Set(ctx, key, 1, constants.TOTPFailureReset).Err()
	} else { // key exists, INCR
		err = c.rdb.Incr(ctx, key).Err()
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *UserCache) TOTPFailureCount(ctx context.Context, uid uint) (int, error) {
	key := c.totpFailureKey(uid)
	ok, err := c.rdb.Exists(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if ok == 0 { // key not exists, return 0
		return 0, nil
	}
	// key exists, return value of key
	str, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	count, _ := strconv.ParseInt(str, 10, 64)
	return int(count), nil
}

func (c *UserCache) totpPendingKey(uid uint) string {
	return fmt.Sprintf("totp:pending:%d", uid)
}

func (c *UserCache) totpPeriodKey(uid uint, period int64) string {
	return fmt.Sprintf("totp:replay:%d:%d", uid, period)
}

func (c *UserCache) totpFailureKey(uid uint) string {
	return fmt.Sprintf("totp:failure:%d", uid)
}

type UserCache struct {
	rdb *redis.Client
}

func NewUserCache() *UserCache {
	return &UserCache{
		rdb: rd.Rdb,
	}
}
