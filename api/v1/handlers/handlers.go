package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"net/http"
	"steg/api/v1/models"
	files "steg/files"
)

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

// Create a new job for steg, read multi-part file upload
func (h *EncodeHandler) EncodePostHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new job
	job := models.NewEncodeJob()
	// Add job to state
	h.jobs[job.Uuid] = job

	// Read file from multipart form, and save it
	saver := files.NewMultiPartFileSaver(r, "image", "raw", job.Uuid.String(), h.validator)
	name, err := saver.SaveFile()
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	job.RawImagePath = name
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusAccepted)
	// Starts a process in goroutine
	job.ProcessImage()
}

// Get steg state/image
func (h *EncodeHandler) EncodeGetHandler(w http.ResponseWriter, r *http.Request) {
	uuidString := r.URL.Query().Get("uuid")
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
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	t, err := json.MarshalIndent(job, "", "  ")
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write(t)
}

type DecodeHandler struct {
	jobs      map[uuid.UUID]*models.DecodeJob
	validator files.FileValidator
}

func NewDecodeHandler(validator files.FileValidator) *DecodeHandler {
	return &DecodeHandler{jobs: make(map[uuid.UUID]*models.DecodeJob), validator: validator}
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

// Create a new job to decode
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
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusAccepted)
	// Starts a process in goroutine
	job.ProcessImage()
}

// Get decode job/message
func (h *DecodeHandler) DecodeGetHandler(w http.ResponseWriter, r *http.Request) {

}

// Get some generic stats
func GetStats(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, Stats")
}
