package handlers

import (
	"errors"
	"log"
	"net/http"
	"steg/shared/repositories"
	"steg/shared/services"

	"github.com/google/uuid"
)

type JobStartHandler struct {
	jobService    services.JobService
	rabbitService *services.RabbitService
}

func NewJobStartHandler(jobService services.JobService, rabbitService *services.RabbitService) *JobStartHandler {
	return &JobStartHandler{jobService, rabbitService}
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

	// Check that UUID is valid
	exists, err := handler.jobService.Exists(r.Context(), jobUUID)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Println("Error Finding job", err)
		return
	}
	if !exists {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	// Kick off job
	jobType := handler.jobService.GetType()
	err = handler.rabbitService.Publish(r.Context(), jobType, jobUUID)
	if err != nil && !errors.Is(err, repositories.AlreadyInProgress) {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Println("Error pushing to rabbit:", err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		log.Println(err)
	}
}
