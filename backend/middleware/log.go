package middleware

import (
	"log"
	"net/http"
)

// Log all requests
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}
