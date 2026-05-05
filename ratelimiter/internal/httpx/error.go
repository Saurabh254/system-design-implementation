package httpx

import "net/http"

type HTTPError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, HTTPError{Status: status, Message: msg})
}
