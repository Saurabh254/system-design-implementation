package tokenbucket

import "fmt"

// configKey returns the Redis hash key that stores
// capacity and fill_rate for an entity type.
func configKey(entity string) string {
	return fmt.Sprintf("rl:%s:config", entity)
}

// tokensKey returns the Redis key that stores the
// current token count for a specific entity instance.
func tokensKey(entity, entityID string) string {
	return fmt.Sprintf("rl:%s:%s:tokens", entity, entityID)
}

// lastFilledKey returns the Redis key that stores the
// Unix timestamp of the last token refill.
func lastFilledKey(entity, entityID string) string {
	return fmt.Sprintf("rl:%s:%s:last_filled", entity, entityID)
}
