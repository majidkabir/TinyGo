package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logger logs HTTP requests
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log request
		log.Printf("→ %s %s", r.Method, r.URL.Path)

		// Call next handler
		next(w, r)

		// Log response time
		duration := time.Since(start)
		log.Printf("← %s %s completed in %v", r.Method, r.URL.Path, duration)
	}
}

// CORS adds CORS headers
func CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
