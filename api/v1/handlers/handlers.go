package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"steg/api/v1/models"
	"steg/middleware"
)

type StatsHandler struct {
	stats *middleware.Stats
	ds    *models.Datastore
}

func NewStatsHandler(stats *middleware.Stats, ds *models.Datastore) *StatsHandler {
	return &StatsHandler{stats, ds}
}

func (h *StatsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if (*r).Method != "GET" {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// TODO Break down jobs by status
	stats := map[string]interface{}{
		"datastore": map[string]interface{}{
			"embed-jobs":   len(h.ds.EmbedJobs),
			"extract-jobs": len(h.ds.ExtractJobs),
			"images":       len(h.ds.Images),
		},
		"requests": h.stats.GetStats(),
	}
	err := json.NewEncoder(w).Encode(stats)
	if err != nil {
		fmt.Println(err)
	}
}
