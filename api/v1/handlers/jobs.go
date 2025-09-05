package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"steg/api/v1/services"

	"github.com/google/uuid"
)

type JobStartHandler[T any] struct {
	service *services.JobService[T]
}

func NewJobStartHandler[T any](service *services.JobService[T]) *JobStartHandler[T] {
	return &JobStartHandler[T]{service: service}
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

	// Kicks off the job
	job, err := j.service.Start(jobUUID)
	if job == nil || errors.Is(err, NotFoundError) {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(job)
	if err != nil {
		fmt.Println(err)
	}
}
