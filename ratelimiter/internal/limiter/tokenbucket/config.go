package tokenbucket

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

const (
	DefaultCapacity = 10
	DefaultFillRate = 1
)

// Config holds the rate limit parameters for an entity type.
// It is stored as a Redis hash and shared across all IDs of that entity.
type Config struct {
	Capacity int `json:"capacity"`
	FillRate int `json:"fill_rate"`
}

// LoadConfig reads the rate limit config for an entity type from Redis.
// If no config has been saved yet, the package defaults are returned.
func LoadConfig(ctx context.Context, rdb *redis.Client, entity string) (Config, error) {
	vals, err := rdb.HGetAll(ctx, configKey(entity)).Result()
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		Capacity: DefaultCapacity,
		FillRate: DefaultFillRate,
	}

	if v, ok := vals["capacity"]; ok {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.Capacity = n
		}
	}
	if v, ok := vals["fill_rate"]; ok {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.FillRate = n
		}
	}

	return cfg, nil
}

// SaveConfig writes the rate limit config for an entity type to Redis.
// The change takes effect on the very next Consume call — no restart needed.
func SaveConfig(ctx context.Context, rdb *redis.Client, entity string, cfg Config) error {
	return rdb.HSet(ctx, configKey(entity),
		"capacity", cfg.Capacity,
		"fill_rate", cfg.FillRate,
	).Err()
}
