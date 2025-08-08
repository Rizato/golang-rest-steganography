package crud

import (
	"github.com/google/uuid"
	"io"
)

// TODO break this up into smaller parts that are composable within an overall interface (CRUDRouter)

// GenericCRUD interface that implements crud for a given resource
type GenericCRUD[T any] interface {
	Create(io.Reader) (T, error)        // POST
	Read(id uuid.UUID) (T, bool, error) // GET
	Update(id uuid.UUID) (T, error)     // PUT (Should replace the full object, so I don't like unless doing PUT /obj/field/
	Delete(id uuid.UUID) (bool, error)  // DELETE
}
