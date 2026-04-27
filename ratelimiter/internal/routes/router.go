package routes

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/saurabh254/system-design-implementation/ratelimiter/internal/middleware"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/docs/", httpSwagger.WrapHandler)
	mux.HandleFunc("/health", healthHandler)

	handler := middleware.Chain(
		mux,
		middleware.EnrichHeaders,
		middleware.LogRequest,
	)

	return handler
}
