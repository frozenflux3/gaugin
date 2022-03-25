package rest

import "net/http"

func defaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'unsafe-inline';")
}