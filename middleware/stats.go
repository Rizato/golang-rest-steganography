package middleware

import (
	"net/http"
	"time"
)

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
