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

func writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// TODO Move this into a package with GenericCrud
type GenericCrudHandler[T models.Model] struct {
	crud.GenericCRUD[T]
}

func NewGenericJobHandler[T models.Model](crud crud.GenericCRUD[T]) *GenericCrudHandler[T] {
	return &GenericCrudHandler[T]{crud}
}

func (h *GenericCrudHandler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if (*r).Method == "GET" {
		// TODO List handling?
		h.GetHandler(w, r)
	} else if (*r).Method == "POST" {
		h.PostHandler(w, r)
	} else if (*r).Method == "PUT" {
		h.PutHandler(w, r)
	} else if (*r).Method == "DELETE" {
		h.DeleteHandler(w, r)
	} else {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *GenericCrudHandler[T]) GetHandler(w http.ResponseWriter, r *http.Request) {
	jobId := r.PathValue("id")
	jobUUID, err := uuid.Parse(jobId)
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
	}
	job, found, err := h.Read(jobUUID)
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	if !found {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	err = writeJSON(w, http.StatusAccepted, job)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *GenericCrudHandler[T]) PostHandler(w http.ResponseWriter, r *http.Request) {
	// Decode json into proper job request struct

	// Validate request, including check that the file exists

	// Create a new job, grabbing relevant parts of the request struct
	job, err := h.Create(r.Body)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = writeJSON(w, http.StatusCreated, job)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (h *GenericCrudHandler[T]) PutHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
}

func (h *GenericCrudHandler[T]) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	jobId := r.PathValue("id")
	jobUUID, err := uuid.Parse(jobId)
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
	}

	found, err := h.Delete(jobUUID)
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	if !found {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	_, err = w.Write([]byte{})

	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func NewEmbedCrudHandler(ds *models.Datastore) *GenericCrudHandler[*models.EmbedJob] {
	return NewGenericJobHandler[*models.EmbedJob](NewEmbedJobCrud(ds))
}

func NewExtractCrudHandler(ds *models.Datastore) *GenericCrudHandler[*models.ExtractJob] {
	return NewGenericJobHandler[*models.ExtractJob](NewExtractJobCrud(ds))
}

func NewFileCrudHandler(ds *models.Datastore) *GenericCrudHandler[*models.ServerFile] {
	return NewGenericJobHandler[*models.ServerFile](NewFileCrud(ds))
}

// TODO Also dump job/file counts
type StatsHandler struct {
	stats *middleware.Stats
}

func NewStatsHandler(stats *middleware.Stats) *StatsHandler {
	return &StatsHandler{stats}
}

func (h *StatsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if (*r).Method != "GET" {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	}
	err := writeJSON(w, http.StatusOK, h.stats.GetStats())
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	}
}
