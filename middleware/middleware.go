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

// Apply Create a nested Handler of Middleware around the given original handler
func Apply(wrapped http.Handler, extra ...Middleware) http.Handler {
	var h = wrapped
	// We wrap each layer in more and more handlers, with are then run outside in
	for _, middleware := range DefaultMiddleware {
		h = middleware.Wrap(h)
	}
	for _, extraMiddleware := range extra {
		h = extraMiddleware.Wrap(h)
	}
	return h
}

// Apply Create a nested Handler of Middleware around the given original handler func
func ApplyFunc(wrapped func(w http.ResponseWriter, req *http.Request), extra ...Middleware) http.Handler {
	return Apply(http.HandlerFunc(wrapped), extra...)
}

// DefaultMiddleware The default Middleware declaration
var DefaultMiddleware = []Middleware{NewStatsMiddleware(), MiddlewareFunc(LoggerMiddleware)}
