package schemas

import "errors"

type RateLimitRequest struct {
	EntityID   string `json:"entity_id"`
	EntityType string `json:"entity_type"`
}

func (r RateLimitRequest) Validate() error {
	if r.EntityID == "" {
		return errors.New("entity_id is required")
	}

	if r.EntityType == "" {
		return errors.New("entity_type is required")
	}

	return nil
}

type RateLimitResponse struct {
	IsAllowed         bool `json:"is_allowed"`
	TokensRemaining   int  `json:"tokens_remaining"`
	TokensCapacity    int  `json:"tokens_capacity"`
	RetryAfterSeconds int  `json:"retry_after_seconds"`
}
