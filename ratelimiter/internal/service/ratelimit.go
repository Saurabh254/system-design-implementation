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
