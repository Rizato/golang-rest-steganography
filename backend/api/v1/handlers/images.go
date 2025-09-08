package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"steg/api/v1/models"
	"steg/api/v1/services"
	files "steg/files"

	"github.com/google/uuid"
)

// FileCrud implements RD of crud, because we have separate handlers for file content
type FileCrud struct {
	imageService *services.ImageService
}

func NewFileCrud(imageService *services.ImageService) *FileCrud {
	return &FileCrud{imageService}
}

func (h *FileCrud) Read(ctx context.Context, uuid uuid.UUID) (*models.Image, error) {
	image, err := h.imageService.GetImage(ctx, uuid)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, NotFoundError
	}
	return image, err
}

func (h *FileCrud) Delete(ctx context.Context, uuid uuid.UUID) error {
	err := h.imageService.DeleteImage(ctx, uuid)
	if errors.Is(err, sql.ErrNoRows) {
		return NotFoundError
	}
	return err
}

// HandleUpload Handles user file uploads from a multipart form, adds a new ImageFile to the application state
func HandleUpload(service *services.ImageService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Get the uploaded from the form
		uploaded, uploadedHeader, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer uploaded.Close()

		maxSize := files.DefaultValidator.MaxFilesize
		// Use LimitReader to stop reading beyond file size, + 1 so I can detect oversized files with spoofed values
		limited := io.LimitReader(uploaded, maxSize+1)

		// Validate the size and mime types (though user supplied and can be manipulated)
		err = files.Validate(uploadedHeader)
		if err != nil {
			if errors.Is(err, files.FileTooLargeError) {
				http.Error(w, "413 Content Too Large", http.StatusRequestEntityTooLarge)
				return
			}

			if errors.Is(err, files.InvalidMimetype) {
				http.Error(w, "400 Bad Request", http.StatusBadRequest)
				return
			}
		}

		// Write to a temporary file
		file, err := os.CreateTemp("", "image")
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			log.Println("Failed to create temporary file", err)
			return
		}
		defer file.Close()
		fileSize, err := io.Copy(file, limited)
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			log.Println("Failed to write temporary file", err)
			return
		}

		// Don't trust client supplied size
		if fileSize > maxSize {
			http.Error(w, "413 File too large", http.StatusRequestEntityTooLarge)
			return
		}

		// Don't trust user supplied mimetype
		_, err = file.Seek(0, 0)
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			log.Println("Error while seeking", err)
			return
		}

		mimetype, err := files.GetMimetype(file)
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			log.Println("Error while determining mimetype", err)
			return
		}

		err = files.ValidateMimetype(mimetype)
		if err != nil {
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
			return
		}

		// Add file to the database
		image, err := service.AddImage(r.Context(), file.Name(), mimetype, fileSize, false)
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			log.Println("Failed to add image", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(image)
		if err != nil {
			fmt.Println(err)
		}
	})
}

func HandleDownload(service *services.ImageService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get image from datastore by id
		fileId := r.PathValue("id")
		fileUUID, err := uuid.Parse(fileId)
		if err != nil {
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
			return
		}

		image, err := service.GetImage(r.Context(), fileUUID)
		if errors.Is(err, NotFoundError) {
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			log.Println("Failed to get image", err)
			return
		}

		// Open file
		f, err := os.Open(image.Path)
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			log.Println("Failed to open file", err)
			return
		}
		defer f.Close()

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", image.Mimetype)
		w.Header().Set("Cache-Control", "public, max-age=86400") // Cache for 1 day

		_, err = io.Copy(w, f)
		if err != nil {
			fmt.Println(err)
		}
	})
}
