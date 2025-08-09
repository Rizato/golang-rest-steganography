package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"steg/api/v1/models"
	"steg/crud"
	"steg/middleware"
)

type JobStartHandler[T models.Job] struct {
	reader crud.Reader[T]
}

func NewJobStartHandler[T models.Job](reader crud.Reader[T]) *JobStartHandler[T] {
	return &JobStartHandler[T]{reader}
}

func (j *JobStartHandler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if (*r).Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	}
	jobId := r.PathValue("id")
	jobUUID, err := uuid.Parse(jobId)
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	job, found, err := j.reader.Read(jobUUID)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !found {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	// Kicks off the job
	job.Start()

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(job)
	if err != nil {
		fmt.Println(err)
	}
}

type StatsHandler struct {
	stats *middleware.Stats
	ds    *models.Datastore
}

func NewStatsHandler(stats *middleware.Stats, ds *models.Datastore) *StatsHandler {
	return &StatsHandler{stats, ds}
}

func (h *StatsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if (*r).Method != http.MethodGet {
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
