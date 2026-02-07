package user

import (
	"StreamCore/internal/pkg/constants"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserCache interface {
	SetTOTPPending(ctx context.Context, uid uint, secret string) error
	GetTOTPPending(ctx context.Context, uid uint) (string, error)
	LockCurrentTOTPPeriod(ctx context.Context, uid uint) error
	IsCurrentTOTPPeriodLocked(ctx context.Context, uid uint) bool
	IncreaseTOTPFailure(ctx context.Context, uid uint) error
	TOTPFailureCount(ctx context.Context, uid uint) (int, error)
}

func NewUserCache(rdb *redis.Client) UserCache {
	return &usercache{
		rdb: rdb,
	}
}

func (c *usercache) SetTOTPPending(ctx context.Context, uid uint, secret string) error {
	cacheKey := c.totpPendingKey(uid)
	return c.rdb.Set(ctx, cacheKey, secret, constants.TOTPSecretExpiry).Err()
}

func (c *usercache) GetTOTPPending(ctx context.Context, uid uint) (string, error) {
	cacheKey := c.totpPendingKey(uid)
	return c.rdb.Get(ctx, cacheKey).Result()
}

func (c *usercache) LockCurrentTOTPPeriod(ctx context.Context, uid uint) error {
	timestep := time.Now().UnixMilli() / constants.TOTPInterval
	key := c.totpPeriodKey(uid, timestep)
	return c.rdb.Set(ctx, key, 1, constants.TOTPInterval*2*time.Second).Err()
}

func (c *usercache) IsCurrentTOTPPeriodLocked(ctx context.Context, uid uint) bool {
	timestep := time.Now().UnixMilli() / constants.TOTPInterval
	key := c.totpPeriodKey(uid, timestep)
	ok, _ := c.rdb.Exists(ctx, key).Result()
	return ok == 1
}

func (c *usercache) IncreaseTOTPFailure(ctx context.Context, uid uint) error {
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

func (c *usercache) TOTPFailureCount(ctx context.Context, uid uint) (int, error) {
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

func (c *usercache) totpPendingKey(uid uint) string {
	return fmt.Sprintf("totp:pending:%d", uid)
}

func (c *usercache) totpPeriodKey(uid uint, period int64) string {
	return fmt.Sprintf("totp:replay:%d:%d", uid, period)
}

func (c *usercache) totpFailureKey(uid uint) string {
	return fmt.Sprintf("totp:failure:%d", uid)
}

type usercache struct {
	rdb *redis.Client
}
