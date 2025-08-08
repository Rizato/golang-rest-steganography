package crud

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"steg/api/v1/models"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// TODO break this up into smaller parts that are composable within an overall interface (CRUDRouter)

// GenericCRUD interface that implements crud for a given resource
type GenericCRUD[T any] interface {
	Create(io.Reader) (T, error)        // POST
	Read(id uuid.UUID) (T, bool, error) // GET
	Update(id uuid.UUID) (T, error)     // PUT (Should replace the full object, so I don't like unless doing PUT /obj/field/
	Delete(id uuid.UUID) (bool, error)  // DELETE
}

type GenericCrudHandler[T models.Model] struct {
	GenericCRUD[T]
}

func NewGenericJobHandler[T models.Model](crud GenericCRUD[T]) *GenericCrudHandler[T] {
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
