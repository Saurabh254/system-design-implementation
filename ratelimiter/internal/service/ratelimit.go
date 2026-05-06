package service

import (
	"context"
	"net/http"

	"github.com/saurabh254/system-design-implementation/ratelimiter/internal/httpx"
	"github.com/saurabh254/system-design-implementation/ratelimiter/internal/limiter/tokenbucket"
	"github.com/saurabh254/system-design-implementation/ratelimiter/internal/schemas"
	store "github.com/saurabh254/system-design-implementation/ratelimiter/internal/store/redis"
)

func GetRateLimitStatus(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req schemas.RateLimitRequest

	req.EntityID = r.PathValue("entity_id")
	req.EntityType = r.PathValue("entity_type")

	if req.EntityID == "" || req.EntityType == "" {
		httpx.Error(w, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	rdb := store.NewClient()
	tb := tokenbucket.New(req.EntityType, req.EntityID, rdb)

	response, err := tb.Status(ctx)
	if err != nil {
		http.Error(w, "failed to get rate limit status", http.StatusInternalServerError)
		return
	}
	httpx.JSON(w, http.StatusOK, response)
}

func ConsumeRateLimit(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req schemas.RateLimitRequest

	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.EntityID == "" || req.EntityType == "" {
		httpx.Error(w, http.StatusBadRequest, "Entity ID and type are required")
		return
	}

	rdb := store.NewClient()
	tb := tokenbucket.New(req.EntityType, req.EntityID, rdb)

	response, err := tb.Consume(ctx)
	if err != nil {
		http.Error(w, "failed to consume rate limit", http.StatusInternalServerError)
		return
	}

	if response.IsAllowed {
		httpx.JSON(w, http.StatusOK, response)
	} else {
		httpx.JSON(w, http.StatusTooManyRequests, response)
	}
}
