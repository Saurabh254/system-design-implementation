package store

import (
	"github.com/redis/go-redis/v9"
	"github.com/saurabh254/system-design-implementation/ratelimiter/internal/config"
)

type RedisStore struct {
	Client *redis.Client
}

var redisClient *redis.Client

func init() {
	redisClient = NewClient()
}

func NewClient() *redis.Client {

	if redisClient != nil {
		return redisClient
	}
	config := config.Load()
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: "",
		DB:       0,
	})
	return rdb
}
