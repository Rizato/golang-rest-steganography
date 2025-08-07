package middleware

import (
	"net/http"
	"sync"
	"time"
)

type Stats struct {
	mutex         sync.RWMutex
	RequestCount  map[string]int64
	ResponseTime  map[string]time.Duration
	RequestsTotal int64
	StartTime     time.Time
}

func NewStats() *Stats {
	return &Stats{RequestCount: make(map[string]int64), ResponseTime: make(map[string]time.Duration), StartTime: time.Now()}
}

func (s *Stats) TrackCall(method string, path string, duration time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.RequestCount[method+path]++
	s.ResponseTime[method+path] += duration
}

func (s *Stats) GetStats() map[string]interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	stats := map[string]interface{}{
		"uptime":         time.Since(s.StartTime),
		"total_requests": s.RequestsTotal,
		"endpoints":      make(map[string]interface{}),
	}

	endpoints := stats["endpoints"].(map[string]interface{})
	for endpoint, count := range s.RequestCount {
		avgDuration := time.Duration(0)
		if count > 0 {
			avgDuration = s.ResponseTime[endpoint] / time.Duration(count)
		}

		endpoints[endpoint] = map[string]interface{}{
			"count":      count,
			"total_time": s.ResponseTime[endpoint].String(),
			"avg_time":   avgDuration.String(),
		}
	}

	return stats
}

type StatsMiddleware struct {
	stats *Stats
}

func NewStatsMiddleware(s *Stats) *StatsMiddleware {
	return &StatsMiddleware{stats: s}
}

func (s StatsMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		next.ServeHTTP(w, r)
		end := time.Now()
		s.stats.TrackCall(r.Method, r.URL.Path, end.Sub(start))
	})
}
