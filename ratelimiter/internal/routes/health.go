package routes

import "net/http"

// HealthCheck godoc
// @Summary Health check
// @Description check server status
// @Tags health
// @Success 200 {string} string "ok"
// @Router /health [get]
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
