package main

import (
	"log"
	"net/http"
	"time"
)

func LogMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Call the next handler in the chain
		next.ServeHTTP(w, r)

		duration := time.Since(start)
		log.Printf("Request %s %s took %v", r.Method, r.URL.String(), duration)
	})
}
