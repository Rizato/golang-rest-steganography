package crud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

type Creator[T any] interface {
	Create(ctx context.Context, reader io.Reader) (T, error) // POST
}

type Reader[T any] interface {
	Read(ctx context.Context, id uuid.UUID) (T, bool, error) // GET
}

type Updater[T any] interface {
	Update(ctx context.Context, id uuid.UUID, reader io.Reader) (T, bool, error) // PUT
}

type Deleter[T any] interface {
	Delete(ctx context.Context, id uuid.UUID) (bool, error) // DELETE
}

type Lister[T any] interface {
	List(ctx context.Context) ([]T, error)
}

type BaseCrud[T any] interface {
}

type ListHandler[T any] struct {
	service BaseCrud[T]
}

func NewListHandler[T any](crud BaseCrud[T]) *ListHandler[T] {
	return &ListHandler[T]{crud}
}

func (h *ListHandler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// TODO Handle Options
	case http.MethodPost:
		if creator, ok := any(h.service).(Creator[T]); ok {
			h.PostHandler(w, r, creator)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case http.MethodGet:
		if lister, ok := any(h.service).(Lister[T]); ok {
			h.GetHandler(w, r, lister)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		break
	}
}

func (h *ListHandler[T]) PostHandler(w http.ResponseWriter, r *http.Request, creator Creator[T]) {
	resource, err := creator.Create(r.Context(), r.Body)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = writeJSON(w, http.StatusCreated, resource)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (h *ListHandler[T]) GetHandler(w http.ResponseWriter, r *http.Request, lister Lister[T]) {
	resources, err := lister.List(r.Context())
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	err = writeJSON(w, http.StatusAccepted, resources)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
}

type ItemHandler[T any] struct {
	service BaseCrud[T]
}

func NewItemHandler[T any](crud BaseCrud[T]) *ItemHandler[T] {
	return &ItemHandler[T]{crud}
}

func (h *ItemHandler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// TODO options
	case http.MethodGet:
		if reader, ok := any(h.service).(Reader[T]); ok {
			h.GetHandler(w, r, reader)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case http.MethodPut:
		if updater, ok := any(h.service).(Updater[T]); ok {
			h.PutHandler(w, r, updater)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case http.MethodDelete:
		if deleter, ok := any(h.service).(Deleter[T]); ok {
			h.DeleteHandler(w, r, deleter)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		break
	}
}

func (h *ItemHandler[T]) GetHandler(w http.ResponseWriter, r *http.Request, reader Reader[T]) {
	resourceId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
	}
	resource, found, err := reader.Read(r.Context(), resourceId)
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	if !found {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	err = writeJSON(w, http.StatusAccepted, resource)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *ItemHandler[T]) PutHandler(w http.ResponseWriter, r *http.Request, updater Updater[T]) {
	resourceUUID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
	}

	resource, found, err := updater.Update(r.Context(), resourceUUID, r.Body)
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	if !found {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	err = writeJSON(w, http.StatusCreated, resource)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (h *ItemHandler[T]) DeleteHandler(w http.ResponseWriter, r *http.Request, deleter Deleter[T]) {
	resourceUUID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
	}

	found, err := deleter.Delete(r.Context(), resourceUUID)
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
