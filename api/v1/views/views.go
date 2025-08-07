package views

import (
	"net/http"
	"steg/api/v1/handlers"
	files "steg/files"
	"steg/middleware"
)

func ConfigureViews(stats *middleware.Stats) http.Handler {
	mux := http.NewServeMux()
	validator := files.NewFileValidator(files.DefaultMaxFilesize, files.DefaultMimeTypes...)

	encodeHandler := handlers.NewEncodeHandler(validator)
	// POST create a new job to steg the given image
	mux.Handle("POST /api/v1/encode", encodeHandler)

	// Get details about a job, and possibly the stegg'd jpeg
	mux.Handle("GET /api/v1/encode/{id}", encodeHandler)

	// Post create a job to read the encoded message, if found
	decodeHandler := handlers.NewDecodeHandler(validator)
	mux.Handle("POST /api/v1/decode", decodeHandler)

	// Get the message if one was stegg'd
	mux.Handle("GET /api/v1/decode/{id}", decodeHandler)

	// Liveness check, and stats
	statsHandler := handlers.NewStatsHandler(stats)
	mux.Handle("GET /api/v1/stats", statsHandler)
	return mux
}
