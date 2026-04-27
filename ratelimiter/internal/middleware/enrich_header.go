package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type ctxKey string

const (
	requestIDKey ctxKey = "request_id"
	startTimeKey ctxKey = "start_time"
)

func EnrichHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		start := time.Now()

		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		ctx = context.WithValue(ctx, startTimeKey, start)

		// propagate outward
		w.Header().Set("X-Request-ID", reqID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
