package service

import (
	"net/http"

	"github.com/saurabh254/system-design-implementation/ratelimiter/internal/httpx"
)

func HealthService(w http.ResponseWriter, r *http.Request) {
	httpx.JSON(w, http.StatusAccepted, map[string]string{"status": "ok"})

}
