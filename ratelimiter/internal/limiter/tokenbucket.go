package tokenbucket

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// Maximum number of tokens a bucket can hold.
	Capacity = 10

	// Number of tokens added per second.
	FillRate = 1
)

// TokenBucket represents a rate limiter bucket for a single user.
type TokenBucket struct {
	capacity int
	tokens   int
	fillRate int

	// Unix timestamp of the last refill.
	lastFilled int64

	userID string
	rdb    *redis.Client
}

// GetTokenBucket loads a user's token bucket from Redis.
// If the bucket does not exist, default values are initialized.
func GetTokenBucket(
	ctx context.Context,
	userID string,
	rdb *redis.Client,
) (*TokenBucket, error) {
	tokens, err := getTokens(ctx, rdb, userID)
	if err != nil {
		return nil, err
	}

	lastFilled, err := getLastFilled(ctx, rdb, userID)
	if err != nil {
		return nil, err
	}

	return &TokenBucket{
		capacity:   Capacity,
		tokens:     tokens,
		fillRate:   FillRate,
		lastFilled: lastFilled,
		userID:     userID,
		rdb:        rdb,
	}, nil
}

// getTokens fetches the current token count from Redis.
// If the key does not exist, the bucket is initialized with full capacity.
func getTokens(
	ctx context.Context,
	rdb *redis.Client,
	userID string,
) (int, error) {
	key := userTokensKey(userID)

	value, err := rdb.Get(ctx, key).Result()

	switch {
	case err == redis.Nil:
		err = rdb.Set(ctx, key, Capacity, 0).Err()
		if err != nil {
			return 0, err
		}

		return Capacity, nil

	case err != nil:
		return 0, err
	}

	tokens, err := strconv.Atoi(value)

	// If Redis contains invalid data,
	// fallback to default capacity.
	if err != nil {
		return Capacity, nil
	}

	return tokens, nil
}

// getLastFilled fetches the last refill timestamp from Redis.
// If the key does not exist, current time is stored and returned.
func getLastFilled(
	ctx context.Context,
	rdb *redis.Client,
	userID string,
) (int64, error) {
	key := userLastFilledKey(userID)

	value, err := rdb.Get(ctx, key).Result()

	switch {
	case err == redis.Nil:
		now := time.Now().Unix()

		err = rdb.Set(ctx, key, now, 0).Err()
		if err != nil {
			return 0, err
		}

		return now, nil

	case err != nil:
		return 0, err
	}

	lastFilled, err := strconv.ParseInt(value, 10, 64)

	// If Redis contains invalid data,
	// fallback to current time.
	if err != nil {
		return time.Now().Unix(), nil
	}

	return lastFilled, nil
}

// IsAllowed checks whether a request can consume a token.
// Returns true if the request is allowed.
func (tb *TokenBucket) IsAllowed() bool {
	// Refill tokens before checking availability.
	tb.refresh()

	// Reject request if bucket is empty.
	if tb.tokens <= 0 {
		return false
	}

	// Consume one token for the request.
	tb.consume(1)

	return true
}

// refresh adds tokens based on elapsed time since last refill.
func (tb *TokenBucket) refresh() {
	now := time.Now().Unix()

	// Calculate elapsed time in seconds.
	elapsed := now - tb.lastFilled

	// Determine how many tokens should be added.
	tokensToAdd := int(elapsed) * tb.fillRate

	// Skip update if no refill is needed.
	if tokensToAdd <= 0 {
		return
	}

	// Prevent token count from exceeding capacity.
	newTokenCount := min(
		tb.capacity,
		tb.tokens+tokensToAdd,
	)

	tb.setTokens(newTokenCount)
	tb.setLastFilled(now)
}

// consume removes tokens from the bucket.
func (tb *TokenBucket) consume(count int) {
	tb.setTokens(tb.tokens - count)
}

// setTokens updates the token count
// both locally and in Redis.
func (tb *TokenBucket) setTokens(tokens int) {
	tb.tokens = tokens

	tb.rdb.Set(
		context.Background(),
		userTokensKey(tb.userID),
		tokens,
		0,
	)
}

// setLastFilled updates the refill timestamp
// both locally and in Redis.
func (tb *TokenBucket) setLastFilled(lastFilled int64) {
	tb.lastFilled = lastFilled

	tb.rdb.Set(
		context.Background(),
		userLastFilledKey(tb.userID),
		lastFilled,
		0,
	)
}

// userTokensKey returns the Redis key
// used for storing token count.
func userTokensKey(userID string) string {
	return "user:" + userID + ":tokens"
}

// userLastFilledKey returns the Redis key
// used for storing last refill timestamp.
func userLastFilledKey(userID string) string {
	return "user:" + userID + ":lastFilled"
}

// min returns the smaller integer value.
func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
