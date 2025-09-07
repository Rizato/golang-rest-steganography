package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"steg/api/v1/repository"
	"steg/api/v1/services"

	"github.com/google/uuid"
)

type JobStartHandler[T any] struct {
	service services.JobService[T]
}

func NewJobStartHandler[T any](service services.JobService[T]) *JobStartHandler[T] {
	return &JobStartHandler[T]{service: service}
}

func (handler *JobStartHandler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if (*r).Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	jobId := r.PathValue("id")
	jobUUID, err := uuid.Parse(jobId)
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	// Kicks off the job
	job, err := handler.service.GetJob(r.Context(), jobUUID)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	job, err = handler.service.Start(r.Context(), job)
	if err != nil && !errors.Is(err, repository.AlreadyInProgress) {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(job)
	if err != nil {
		fmt.Println(err)
	}
}
