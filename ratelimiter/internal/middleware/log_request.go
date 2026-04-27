package middleware

import (
	"net/http"
	"time"

	logger "github.com/saurabh254/system-design-implementation/ratelimiter/internal/utils"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		reqID, _ := r.Context().Value(requestIDKey).(string)
		start, _ := r.Context().Value(startTimeKey).(time.Time)

		var duration time.Duration
		if !start.IsZero() {
			duration = time.Since(start)
		}

		logger.Log.Info("request",
			"request_id", reqID,
			"method", r.Method,
			"path", r.URL.Path,
			"duration", duration.String(),
		)
	})
}
