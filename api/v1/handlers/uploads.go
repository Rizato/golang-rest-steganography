package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"steg/api/v1/models"
	files "steg/files"
)

// HandleUpload Handles user file uploads from a multipart form, adds a new ImageFile to the application state
func HandleUpload(ds *models.Datastore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		}

		// Get the uploaded from the form
		uploaded, uploadedHeader, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
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
			}

			if errors.Is(err, files.InvalidMimetype) {
				http.Error(w, "400 Bad Request", http.StatusBadRequest)
			}
		}

		// Write to an specified dir
		image := models.NewServerFile()
		file, err := os.Create("./uploads/" + image.GetUUID().String())
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		}
		defer file.Close()
		fileSize, err := io.Copy(file, limited)
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		}

		// Don't trust client supplied size
		if fileSize > maxSize {
			http.Error(w, "413 File too large", http.StatusRequestEntityTooLarge)
		}

		// Don't trust user supplied mimetype
		_, err = file.Seek(0, 0)
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		}

		mimetype, err := files.GetMimetype(file)
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		}

		err = files.ValidateMimetype(mimetype)
		if err != nil {
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
		}

		image.Path = file.Name()
		image.Size = fileSize
		image.Mimetype = mimetype
		ds.Images[image.GetUUID()] = image

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(image)
		if err != nil {
			fmt.Println(err)
		}
	})
}

func HandleDownload(ds *models.Datastore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get image from datastore by id
		fileId := r.PathValue("id")
		fileUUID, err := uuid.Parse(fileId)
		if err != nil {
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
		}

		image, found := ds.Images[fileUUID]
		if !found {
			http.Error(w, "404 Not Found", http.StatusNotFound)
		}

		// Open file
		f, err := os.Open(image.Path)
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		}
		defer f.Close()

		// Find mimetype
		mimetype, err := files.GetMimetype(f)
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		}

		// Reset file position
		_, err = f.Seek(0, 0)
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", mimetype)
		w.Header().Set("Cache-Control", "public, max-age=86400") // Cache for 1 day

		_, err = io.Copy(w, f)
		if err != nil {
			fmt.Println(err)
		}
	})
}
