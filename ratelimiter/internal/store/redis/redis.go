package store

import (
	"github.com/redis/go-redis/v9"
	"github.com/saurabh254/system-design-implementation/ratelimiter/internal/config"
)

type RedisStore struct {
	Client *redis.Client
}

func NewClient() *redis.Client {
	config := config.Load()
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: "",
		DB:       0,
	})
	return rdb
}
