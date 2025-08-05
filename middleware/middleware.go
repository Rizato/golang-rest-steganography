package middleware

import (
	"fmt"
	"net/http"
)

type Middleware = func(http.Handler) http.Handler

var DefaultMiddleware = []Middleware{LoggerMiddleware, StatsMiddleware}

func Apply(wrapped func(w http.ResponseWriter, req *http.Request), extra ...Middleware) http.Handler {
	var h http.Handler = http.HandlerFunc(wrapped)
	for _, middleware := range DefaultMiddleware {
		h = middleware(h)
	}
	for _, extraMiddleware := range extra {
		h = extraMiddleware(h)
	}
	return h
}

// Log all requests
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

// Update the stats holder, should probably be pass as an arg though
func StatsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}
