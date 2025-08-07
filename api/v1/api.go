package v1

import (
	"net/http"
	handlers "steg/handlers/v1"
	"steg/middleware"
)

func ConfigureServer(server *http.ServeMux) {
	encodeHandler := handlers.NewEncodeHandler()
	// POST create a new job to steg the given image
	server.Handle("POST /api/v1/steganography", middleware.Apply(encodeHandler))

	// Get details about a job, and possibly the stegg'd jpeg
	server.Handle("GET /api/v1/steganography/{id}", middleware.Apply(encodeHandler))

	// Post create a job to read the encoded message, if found
	server.Handle("POST /api/v1/decode", middleware.ApplyFunc(handlers.PostDecodeHandler))

	// Get the message if one was stegg'd
	server.Handle("GET /api/v1/decode/{id}", middleware.ApplyFunc(handlers.GetDecodeHandler))

	// Liveness check, and stats
	server.Handle("/api/v1/stats", middleware.ApplyFunc(handlers.GetStats))
}
