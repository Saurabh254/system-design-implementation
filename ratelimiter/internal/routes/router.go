package routes

import "net/http"

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	registerHealth(mux)

	return mux
}
