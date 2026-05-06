package tokenbucket

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/saurabh254/system-design-implementation/ratelimiter/internal/schemas"
)

// consumeScript atomically refills, checks, and consumes one token.
//
// KEYS[1] — tokens key
// KEYS[2] — last_filled key
// KEYS[3] — config hash key
//
// ARGV[1] — default capacity  (used when no config key exists in Redis)
// ARGV[2] — default fill_rate (used when no config key exists in Redis)
// ARGV[3] — current Unix timestamp
//
// Returns a four-element array: { allowed, tokens_remaining, capacity, retry_after }
var consumeScript = redis.NewScript(`
local tokens_key      = KEYS[1]
local last_filled_key = KEYS[2]
local config_key      = KEYS[3]

-- Read config from Redis; fall back to caller-supplied defaults.
local capacity  = tonumber(redis.call('HGET', config_key, 'capacity'))  or tonumber(ARGV[1])
local fill_rate = tonumber(redis.call('HGET', config_key, 'fill_rate')) or tonumber(ARGV[2])
local now       = tonumber(ARGV[3])

-- Load bucket state, or seed a fresh bucket at full capacity.
local raw_tokens      = redis.call('GET', tokens_key)
local raw_last_filled = redis.call('GET', last_filled_key)

local tokens, last_filled
if raw_tokens == false then
    tokens      = capacity
    last_filled = now
else
    tokens      = tonumber(raw_tokens)
    last_filled = tonumber(raw_last_filled) or now
end

-- Refill tokens proportional to elapsed time.
local elapsed = now - last_filled
local to_add  = math.floor(elapsed * fill_rate)
if to_add > 0 then
    tokens      = math.min(capacity, tokens + to_add)
    last_filled = now
end

-- Consume one token if available.
if tokens > 0 then
    tokens = tokens - 1
    redis.call('SET', tokens_key, tokens)
    redis.call('SET', last_filled_key, last_filled)
    return {1, tokens, capacity, 0}
end

-- Denied — tell the caller how long to wait for the next token.
local retry_after = fill_rate > 0 and math.ceil(1 / fill_rate) or -1
return {0, tokens, capacity, retry_after}
`)

// statusScript reads the current bucket state without consuming a token.
// Useful for health checks or informational endpoints.
//
// KEYS and ARGV are identical to consumeScript.
// Returns: { allowed, tokens_remaining, capacity, retry_after }
var statusScript = redis.NewScript(`
local tokens_key      = KEYS[1]
local last_filled_key = KEYS[2]
local config_key      = KEYS[3]

local capacity  = tonumber(redis.call('HGET', config_key, 'capacity'))  or tonumber(ARGV[1])
local fill_rate = tonumber(redis.call('HGET', config_key, 'fill_rate')) or tonumber(ARGV[2])
local now       = tonumber(ARGV[3])

local raw_tokens      = redis.call('GET', tokens_key)
local raw_last_filled = redis.call('GET', last_filled_key)

local tokens, last_filled
if raw_tokens == false then
    tokens      = capacity
    last_filled = now
else
    tokens      = tonumber(raw_tokens)
    last_filled = tonumber(raw_last_filled) or now
end

-- Compute pending refill without writing it back.
local elapsed = now - last_filled
local to_add  = math.floor(elapsed * fill_rate)
if to_add > 0 then
    tokens = math.min(capacity, tokens + to_add)
end

local allowed       = tokens > 0 and 1 or 0
local retry_after   = (tokens == 0 and fill_rate > 0) and math.ceil(1 / fill_rate) or 0
return {allowed, tokens, capacity, retry_after}
`)

// TokenBucket is a rate limiter for a single entity instance (e.g. user:42).
// All operations are backed by Redis and safe for concurrent use.
type TokenBucket struct {
	entity   string
	entityID string
	rdb      *redis.Client
}

// New returns a TokenBucket for the given entity and ID.
// It does not touch Redis until Consume is called.
func New(entity, entityID string, rdb *redis.Client) *TokenBucket {
	return &TokenBucket{
		entity:   entity,
		entityID: entityID,
		rdb:      rdb,
	}
}

// Status returns the current bucket state without consuming a token.
// Use this for health checks or informational endpoints where you want
// to know the remaining capacity without affecting it.
func (tb *TokenBucket) Status(ctx context.Context) (schemas.RateLimitResponse, error) {
	keys := []string{
		tokensKey(tb.entity, tb.entityID),
		lastFilledKey(tb.entity, tb.entityID),
		configKey(tb.entity),
	}
	args := []any{DefaultCapacity, DefaultFillRate, time.Now().Unix()}

	result, err := statusScript.Run(ctx, tb.rdb, keys, args...).Int64Slice()
	if err != nil {
		return schemas.RateLimitResponse{}, err
	}

	return schemas.RateLimitResponse{
		IsAllowed:         result[0] == 1,
		TokensRemaining:   int(result[1]),
		TokensCapacity:    int(result[2]),
		RetryAfterSeconds: int(result[3]),
	}, nil
}

// Consume atomically refills tokens based on elapsed time, then attempts
// to consume one. It returns the outcome and the current bucket state.
//
// Config is read inside the Lua script on every call, so updates made via
// SaveConfig are reflected immediately without restarting the service.
func (tb *TokenBucket) Consume(ctx context.Context) (schemas.RateLimitResponse, error) {
	keys := []string{
		tokensKey(tb.entity, tb.entityID),
		lastFilledKey(tb.entity, tb.entityID),
		configKey(tb.entity),
	}
	args := []any{DefaultCapacity, DefaultFillRate, time.Now().Unix()}

	result, err := consumeScript.Run(ctx, tb.rdb, keys, args...).Int64Slice()
	if err != nil {
		return schemas.RateLimitResponse{}, err
	}

	return schemas.RateLimitResponse{
		IsAllowed:         result[0] == 1,
		TokensRemaining:   int(result[1]),
		TokensCapacity:    int(result[2]),
		RetryAfterSeconds: int(result[3]),
	}, nil
}
