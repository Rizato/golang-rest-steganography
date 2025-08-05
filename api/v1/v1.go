package v1

import (
	"net/http"
	"steg/handlers"
	"steg/middleware"
)

func ConfigureServer(server *http.ServeMux) {
	// POST create a new job to steg the given image
	server.Handle("POST /api/v1/steganography", middleware.Apply(handlers.PostSteganographyHandler))

	// Get details about a job, and possibly the stegg'd jpeg
	server.Handle("GET /api/v1/steganography/{id}", middleware.Apply(handlers.GetSteganographyHandler))

	// Post create a job to read the encoded message, if found
	server.Handle("POST /api/v1/decode", middleware.Apply(handlers.PostDecodeHandler))

	// Get the message if one was stegg'd
	server.Handle("GET /api/v1/decode/{id}", middleware.Apply(handlers.GetDecodeHandler))

	// Liveness check, and stats
	server.Handle("/api/v1/stats", middleware.Apply(handlers.GetStats))
}
