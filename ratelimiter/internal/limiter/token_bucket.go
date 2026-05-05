package limiter

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	store "github.com/saurabh254/system-design-implementation/ratelimiter/internal/store/redis"
)

const (
	Capacity     = 10
	RefillRate   = 1
	RefillPeriod = 1 // in seconds
)

var ctx = context.Background()

type TokenBucket struct {
	ctx            context.Context
	UserID         string
	rdb            *redis.Client
	lastRefillTime int64
}

func InitTokenBucket(user_id string) *TokenBucket {
	rdb := store.NewClient()

	_, err := rdb.Get(ctx, user_id).Result()
	if err != nil && err.Error() != "redis: nil" {
		return nil
	}
	lasttime_result, err := rdb.Get(ctx, fmt.Sprintf("%s_last_refill", user_id)).Result()
	if err == nil {
		lastRefillTime, err := strconv.ParseInt(lasttime_result, 10, 64)
		if err != nil {
			return nil
		}
		return &TokenBucket{
			ctx:            ctx,
			UserID:         user_id,
			rdb:            rdb,
			lastRefillTime: lastRefillTime,
		}
	}

	if err.Error() == "redis: nil" {

		err := rdb.Set(ctx, user_id, Capacity, 0).Err()
		if err != nil {
			return nil
		}
		err = rdb.Set(ctx, fmt.Sprintf("%s_last_refill", user_id), time.Unix(), 0).Err()
		if err != nil {
			return nil
		}
	}
	return &TokenBucket{
		ctx:    ctx,
		UserID: user_id,
		rdb:    store.NewClient(),
		lastRefillTime: ,
	}
}
