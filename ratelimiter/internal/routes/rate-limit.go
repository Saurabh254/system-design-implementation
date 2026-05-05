package routes

import (
	"net/http"
)

func RateLimitRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/status", getRateLimitHandler)
	// mux.HandleFunc("POST /config", updateRateLimitHandler)
	// mux.HandleFunc("GET /config", getRateLimitConfigHandler)
	// mux.HandleFunc("DELETE /config", deleteRateLimitHandler)

	return mux
}

// ShowAccount godoc
// @Summary      Show an account
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Success      200  {object}  string
// @Router       /api/v1/rate-limit/status [get]
func getRateLimitHandler(w http.ResponseWriter, r *http.Request) {

}

// // UpdateRateLimit godoc
// // @Summary Update rate limit
// // @Description update rate limit configuration
// // @Tags rate-limit
// // @Success 200 {string} string "ok"
// // @Router api/v1/rate-limit/config
// func updateRateLimitHandler(w http.ResponseWriter, r *http.Request) {

// }

// // GetRateLimit godoc
// // @Summary Get rate limit
// // @Description get current rate limit configuration
// // @Tags rate-limit
// // @Success 200 {string} string "ok"
// // @Router api/v1/rate-limit/config
// func getRateLimitConfigHandler(w http.ResponseWriter, r *http.Request) {

// }

// // DeleteRateLimit godoc
// // @Summary Delete rate limit
// // @Description delete rate limit configuration
// // @Tags rate-limit
// // @Success 200 {string} string "ok"
// // @Router /rate-limit/config [get]
// func deleteRateLimitHandler(w http.ResponseWriter, r *http.Request) {
// }
