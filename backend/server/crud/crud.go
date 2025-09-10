package crud

import (
	"context"
	"encoding/json"
	"io"
	"log"
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
	crud BaseCrud[T]
}

func NewListHandler[T any](crud BaseCrud[T]) *ListHandler[T] {
	return &ListHandler[T]{crud}
}

func (h *ListHandler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// TODO Handle Options
	case http.MethodPost:
		creator, ok := h.crud.(Creator[T])
		if ok {
			h.HandlePost(w, r, creator)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case http.MethodGet:
		lister, ok := h.crud.(Lister[T])
		if ok {
			h.HandleList(w, r, lister)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ListHandler[T]) HandlePost(w http.ResponseWriter, r *http.Request, creator Creator[T]) {
	resource, err := creator.Create(r.Context(), r.Body)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Println("Failed to create resource:", err)
		return
	}

	err = writeJSON(w, http.StatusCreated, resource)
	if err != nil {
		log.Println(err)
		return
	}
}

func (h *ListHandler[T]) HandleList(w http.ResponseWriter, r *http.Request, lister Lister[T]) {
	resources, err := lister.List(r.Context())
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	err = writeJSON(w, http.StatusAccepted, resources)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Println("Failed to write response:", err)
		return
	}
}

type ItemHandler[T any] struct {
	crud BaseCrud[T]
}

func NewItemHandler[T any](crud BaseCrud[T]) *ItemHandler[T] {
	return &ItemHandler[T]{crud}
}

func (h *ItemHandler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// TODO options
	case http.MethodGet:
		reader, ok := h.crud.(Reader[T])
		if ok {
			h.HandleGet(w, r, reader)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case http.MethodPut:
		updater, ok := h.crud.(Updater[T])
		if ok {
			h.HandlePut(w, r, updater)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case http.MethodDelete:
		deleter, ok := h.crud.(Deleter[T])
		if ok {
			h.HandleDelete(w, r, deleter)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		break
	}
}

func (h *ItemHandler[T]) HandleGet(w http.ResponseWriter, r *http.Request, reader Reader[T]) {
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

	err = writeJSON(w, http.StatusOK, resource)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Println("Failed to write response:", err)
		return
	}
}

func (h *ItemHandler[T]) HandlePut(w http.ResponseWriter, r *http.Request, updater Updater[T]) {
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
		log.Println(err)
		return
	}
}

func (h *ItemHandler[T]) HandleDelete(w http.ResponseWriter, r *http.Request, deleter Deleter[T]) {
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
		log.Println("Failed to write response:", err)
		return
	}
}
