package utils

import (
	"log"
	"net/http"
	"time"
)

// WithLogging wraps an http.HandlerFunc to log request details
func WithLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// We could wrap ResponseWriter to capture status code, but for now we keep it simple
		next(w, r)
		log.Printf("METHOD=%s PATH=%s REMOTE=%s DURATION=%v", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	}
}
