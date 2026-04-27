package routes

import "net/http"

func registerHealth(
	mux *http.ServeMux,
) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
}
