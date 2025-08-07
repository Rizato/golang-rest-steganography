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

type EncodeHandler struct {
	jobs      map[uuid.UUID]*models.EncodeJob
	validator files.FileValidator
}

func NewEncodeHandler(validator files.FileValidator) *EncodeHandler {
	return &EncodeHandler{jobs: make(map[uuid.UUID]*models.EncodeJob), validator: validator}
}

func (h *EncodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if (*r).Method == "GET" {
		h.EncodeGetHandler(w, r)
	} else if (*r).Method == "POST" {
		h.EncodePostHandler(w, r)
	} else {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *EncodeHandler) CreateJob() *models.EncodeJob {
	job := models.NewEncodeJob()
	// Store job in state
	h.jobs[job.Uuid] = job
	return job
}

// POST /api/v1/encode
// Create a new job for steg, read multi-part file upload
func (h *EncodeHandler) EncodePostHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new job
	job := h.CreateJob()

	// Read file from multipart form, and save it
	saver := files.NewMultiPartFileSaver(r, "image", "raw", job.Uuid.String(), h.validator)
	name, err := saver.SaveFile()
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	job.RawImagePath = name
	err = writeJSON(w, http.StatusAccepted, job)
	if err != nil {
		http.Error(w, "400 Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Starts a process in goroutine
	job.ProcessImage()
}

// GET /api/v1/encode/{id}
func (h *EncodeHandler) EncodeGetHandler(w http.ResponseWriter, r *http.Request) {
	uuidString := r.PathValue("id")
	u, err := uuid.Parse(uuidString)
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	job := h.jobs[u]
	if job == nil {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	err = writeJSON(w, http.StatusAccepted, job)
	if err != nil {
		http.Error(w, "400 Internal Server Error", http.StatusInternalServerError)
		return
	}
}

type DecodeHandler struct {
	jobs      map[uuid.UUID]*models.DecodeJob
	validator files.FileValidator
}

func NewDecodeHandler(validator files.FileValidator) *DecodeHandler {
	return &DecodeHandler{jobs: make(map[uuid.UUID]*models.DecodeJob), validator: validator}
}

func (h *DecodeHandler) CreateJob() *models.DecodeJob {
	job := models.NewDecodeJob()
	// Store job in state
	h.jobs[job.Uuid] = job
	return job
}

func (h *DecodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if (*r).Method == "GET" {
		h.DecodeGetHandler(w, r)
	} else if (*r).Method == "POST" {
		h.DecodePostHandler(w, r)
	} else {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// POST /api/v1/decode
func (h *DecodeHandler) DecodePostHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new job
	job := models.NewDecodeJob()
	// Add job to state
	h.jobs[job.Uuid] = job

	// Read file from multipart form, and save it
	saver := files.NewMultiPartFileSaver(r, "image", "decode", job.Uuid.String(), h.validator)
	name, err := saver.SaveFile()
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	job.ImagePath = name
	err = writeJSON(w, http.StatusAccepted, job)
	if err != nil {
		http.Error(w, "400 Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Starts a process in goroutine
	job.ProcessImage()
}

// GET /api/v1/decode/{id}
func (h *DecodeHandler) DecodeGetHandler(w http.ResponseWriter, r *http.Request) {
	uuidString := r.PathValue("id")
	u, err := uuid.Parse(uuidString)
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	job := h.jobs[u]
	if job == nil {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	err = writeJSON(w, http.StatusAccepted, job)
	if err != nil {
		http.Error(w, "400 Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Get some generic stats
// GET /api/v1/stats
func GetStats(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, Stats")
}
