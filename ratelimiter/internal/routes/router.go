package routes

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/saurabh254/system-design-implementation/ratelimiter/internal/middleware"
)

func v1Router() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/rate-limit/", http.StripPrefix("/rate-limit", RateLimitRouter()))

	return mux
}

func NewRouter() http.Handler {

	// Top-level mux — no prefix required
	mux := http.NewServeMux()
	mux.Handle("/docs/", httpSwagger.WrapHandler)
	mux.Handle("/health", HealthHandlerRouter())
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1Router()))

	return middleware.Chain(
		mux,
		middleware.EnrichHeaders,
		middleware.LogRequest,
	)
}
