package routes

import (
	"net/http"

	"github.com/saurabh254/system-design-implementation/ratelimiter/internal/service"
)

// HealthCheck godoc
// @Summary Health check
// @Description check server status
// @Tags health
// @Success 200 {string} string "ok"
// @Router /health [get]
func healthHandler(w http.ResponseWriter, r *http.Request) {
	service.HealthService(w, r)
}
