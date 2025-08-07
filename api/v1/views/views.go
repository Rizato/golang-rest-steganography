package views

import (
	"net/http"
	"steg/api/v1/handlers"
	files "steg/files"
	"steg/middleware"
)

func ConfigureViews(server *http.ServeMux) {
	validator := files.NewFileValidator(files.DefaultMaxFilesize, files.DefaultMimeTypes...)
	encodeHandler := handlers.NewEncodeHandler(validator)
	// POST create a new job to steg the given image
	server.Handle("POST /api/v1/steganography", middleware.Apply(encodeHandler))

	// Get details about a job, and possibly the stegg'd jpeg
	server.Handle("GET /api/v1/steganography/{id}", middleware.Apply(encodeHandler))

	// Post create a job to read the encoded message, if found
	decodeHandler := handlers.NewDecodeHandler(validator)
	server.Handle("POST /api/v1/decode", middleware.Apply(decodeHandler))

	// Get the message if one was stegg'd
	server.Handle("GET /api/v1/decode/{id}", middleware.Apply(decodeHandler))

	// Liveness check, and stats
	server.Handle("/api/v1/stats", middleware.ApplyFunc(handlers.GetStats))
}
