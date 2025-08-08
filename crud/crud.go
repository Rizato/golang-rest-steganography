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

type Creator[T any] interface {
	Create(io.Reader) (T, error) // POST
}

type Reader[T any] interface {
	Read(uuid.UUID) (T, bool, error) // GET
}

type Updater[T any] interface {
	Update(id uuid.UUID, reader io.Reader) (T, bool, error) // PUT
}

type Deleter[T any] interface {
	Delete(uuid.UUID) (bool, error) // DELETE
}

type Lister[T any] interface {
	List() ([]T, error)
}

// FullCrud interface that implements crud for a given resource
type FullCrud[T any] interface {
	Creator[T]
	Reader[T]
	Updater[T]
	Deleter[T]
	Lister[T]
}

type BaseCrud[T any] interface {
}

type GenericCrudHandler[T models.Model] struct {
	BaseCrud[T]
}

func NewGenericJobHandler[T models.Model](crud BaseCrud[T]) *GenericCrudHandler[T] {
	return &GenericCrudHandler[T]{crud}
}

func (h *GenericCrudHandler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// TODO List and options
	case http.MethodGet:
		h.GetHandler(w, r)
		break
	case http.MethodPost:
		h.PostHandler(w, r)
		break
	case http.MethodPut:
		h.PutHandler(w, r)
		break
	case http.MethodDelete:
		h.DeleteHandler(w, r)
		break
	default:
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		break
	}
}

func (h *GenericCrudHandler[T]) GetHandler(w http.ResponseWriter, r *http.Request) {
	reader, ok := any(h).(Reader[T])
	if !ok {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	}
	jobId := r.PathValue("id")
	jobUUID, err := uuid.Parse(jobId)
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
	}
	job, found, err := reader.Read(jobUUID)
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
	creator, ok := any(h).(Creator[T])
	if !ok {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	}
	job, err := creator.Create(r.Body)
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
	updater, ok := any(h).(Updater[T])
	if !ok {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	}
	jobId := r.PathValue("id")
	jobUUID, err := uuid.Parse(jobId)
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
	}

	job, found, err := updater.Update(jobUUID, r.Body)
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	if !found {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	err = writeJSON(w, http.StatusCreated, job)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (h *GenericCrudHandler[T]) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	deleter, ok := any(h).(Deleter[T])
	if !ok {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	}
	jobId := r.PathValue("id")
	jobUUID, err := uuid.Parse(jobId)
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
	}

	found, err := deleter.Delete(jobUUID)
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
