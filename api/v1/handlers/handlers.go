package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"net/http"
	"steg/api/v1/models"
	files "steg/files"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

type JobFactory[T models.ImageProcessor] interface {
	Create() (T, error)
}

type GenericJobHandler[T models.ImageProcessor] struct {
	jobs         map[uuid.UUID]T
	factory      JobFactory[T]
	formFileName string
	downloadPath string
	validator    files.FileValidator
}

func NewGenericJobHandler[T models.ImageProcessor](jobs map[uuid.UUID]T, factory JobFactory[T], formFileName string, downloadPath string, validator files.FileValidator) *GenericJobHandler[T] {
	return &GenericJobHandler[T]{
		jobs,
		factory,
		formFileName,
		downloadPath,
		validator,
	}
}

func (h *GenericJobHandler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if (*r).Method == "GET" {
		h.GetHandler(w, r)
	} else if (*r).Method == "POST" {
		h.PostHandler(w, r)
	} else {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *GenericJobHandler[T]) GetHandler(w http.ResponseWriter, r *http.Request) {
	jobId := r.PathValue("id")
	job, err := h.GetJob(jobId)
	if err != nil {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	err = writeJSON(w, http.StatusAccepted, job)
	if err != nil {
		http.Error(w, "400 Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *GenericJobHandler[T]) GetJob(jobId string) (T, error) {
	jobUuid, err := uuid.Parse(jobId)
	if err != nil {
		// Get a zero'd value (nil when pointer)
		return h.jobs[uuid.New()], err
	}
	return h.jobs[jobUuid], nil
}

func (h *GenericJobHandler[T]) PostHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new job
	job, err := h.CreateJob()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Read file from multipart form, and save it
	saver := files.NewMultiPartFileSaver(r, h.formFileName, h.downloadPath, job.GetUUID().String(), h.validator)
	name, err := saver.SaveFile()
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	job.SetImagePath(name)
	err = writeJSON(w, http.StatusAccepted, job)
	if err != nil {
		http.Error(w, "400 Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Starts a process in goroutine
	job.ProcessImage()
}

func (h *GenericJobHandler[T]) CreateJob() (T, error) {
	job, err := h.factory.Create()
	if err != nil {
		return job, err
	}

	// Store job in state
	h.jobs[job.GetUUID()] = job
	return job, nil
}

func NewEncodeHandler(validator files.FileValidator) *GenericJobHandler[*models.EncodeJob] {
	return NewGenericJobHandler[*models.EncodeJob](make(map[uuid.UUID]*models.EncodeJob), NewEncodeFactory(), "image", "raw", validator)
}

func NewDecodeHandler(validator files.FileValidator) *GenericJobHandler[*models.DecodeJob] {
	return NewGenericJobHandler[*models.DecodeJob](make(map[uuid.UUID]*models.DecodeJob), NewDecodeFactory(), "image", "decode", validator)
}

// Get some generic stats
// GET /api/v1/stats
func GetStats(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, Stats")
}
