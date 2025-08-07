package middleware

import (
	"net/http"
)

// Middleware Base declaration
type Middleware interface {
	Wrap(http.Handler) http.Handler
}

// MiddlewareFunc Adapter for function Middleware
type MiddlewareFunc func(http.Handler) http.Handler

func (m MiddlewareFunc) Wrap(next http.Handler) http.Handler {
	return m(next)
}

func ConfigureMiddleware(mux http.Handler, stats *Stats) http.Handler {
	server := mux
	server = LoggerMiddleware(server)
	server = NewStatsMiddleware(stats).Wrap(server)
	return server
}
