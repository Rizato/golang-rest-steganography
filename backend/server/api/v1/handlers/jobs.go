package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"steg/shared/repositories"
	"steg/shared/services"

	"github.com/google/uuid"
)

type JobStartHandler struct {
	service services.JobService
}

func NewJobStartHandler(service services.JobService) *JobStartHandler {
	return &JobStartHandler{service: service}
}

func (handler *JobStartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Println("Error Finding job", err)
		return
	}

	// Instead of calling start, push to rabbitmq
	err = handler.service.Execute(r.Context(), job)
	if err != nil && !errors.Is(err, repositories.AlreadyInProgress) {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Println("Error starting job:", err)
		return
	}

	job, err = handler.service.GetJob(r.Context(), jobUUID)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Println("Error Finding job", err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(job)
	if err != nil {
		log.Println(err)
	}
}
