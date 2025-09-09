package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"steg/api/v1/services"
)

type StatsHandler struct {
	stats *services.StatsService
}

func NewStatsHandler(stats *services.StatsService) *StatsHandler {
	return &StatsHandler{stats}
}

func (h *StatsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if (*r).Method != http.MethodGet {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	stats := h.stats.GetStats(r.Context())
	err := json.NewEncoder(w).Encode(stats)
	if err != nil {
		log.Println(err)
	}
}
