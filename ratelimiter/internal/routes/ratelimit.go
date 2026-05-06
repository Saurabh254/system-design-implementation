package routes

import (
	"net/http"

	_ "github.com/saurabh254/system-design-implementation/ratelimiter/internal/schemas"
	"github.com/saurabh254/system-design-implementation/ratelimiter/internal/service"
)

func RateLimitRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/{entity_type}/{entity_id}/status", getRateLimitHandler)
	return mux
}

// RateLimit godoc
//
// @Summary      Show rate limit status
// @Description  Returns the current rate limit status for the specified entity.
// @Tags         ratelimit
// @Accept       json
// @Produce      json
// @Param        entity_type path string true "Entity type"
// @Param        entity_id path string true "Entity ID"
// @Success      200  {object}  schemas.RateLimitResponse
// @Router       /api/v1/ratelimit/{entity_type}/{entity_id}/status [get]
func getRateLimitHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()

	service.GetRateLimitStatus(
		ctx,
		w,
		r,
	)
}
