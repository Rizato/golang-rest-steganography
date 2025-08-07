package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// Base Middleware declaration
type Middleware interface {
	Wrap(http.Handler) http.Handler
}

// Adapter for function Middleware
type MiddlewareFunc func(http.Handler) http.Handler

func (m MiddlewareFunc) Wrap(next http.Handler) http.Handler {
	return m(next)
}

var DefaultMiddleware = []Middleware{NewStatsMiddleware(), MiddlewareFunc(LoggerMiddleware)}

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

func ApplyFunc(wrapped func(w http.ResponseWriter, req *http.Request), extra ...Middleware) http.Handler {
	return Apply(http.HandlerFunc(wrapped), extra...)
}

// Log all requests
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

type StatsMiddleware struct {
	calls               map[string]int64
	averageResponseTime int64
	totalCalls          int64
}

func NewStatsMiddleware() *StatsMiddleware {
	return &StatsMiddleware{}
}

// Update the stats holder, should probably be pass as an arg though
func (s *StatsMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.TrackCall(r.Method, r.URL.Path)
		start := time.Now()
		next.ServeHTTP(w, r)
		end := time.Now()

		duration := end.Sub(start)
		s.UpdateAverageResponseDuration(duration)

	})
}

func (s *StatsMiddleware) TrackCall(method string, path string) {
	s.calls[method+path]++
}

func (s *StatsMiddleware) UpdateAverageResponseDuration(duration time.Duration) {
	var tempDuration = (s.averageResponseTime * s.totalCalls) + duration.Milliseconds()

	// Update rolling average
	s.totalCalls += 1
	s.averageResponseTime = tempDuration / s.totalCalls
}
