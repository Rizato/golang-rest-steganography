package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"steg/middleware"
)

// TODO Also dump job/file counts
type StatsHandler struct {
	stats *middleware.Stats
}

func NewStatsHandler(stats *middleware.Stats) *StatsHandler {
	return &StatsHandler{stats}
}

func (h *StatsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if (*r).Method != "GET" {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(h.stats)
	if err != nil {
		fmt.Println(err)
	}
}
