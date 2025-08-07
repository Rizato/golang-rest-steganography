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

	embedHandler := handlers.NewEmbedHandler(validator)
	// POST create a new job to steg the given image
	mux.Handle("POST /api/v1/embed", http.StripPrefix("/api/v1/", embedHandler))

	// Get details about a job, and possibly the stegg'd jpeg
	mux.Handle("GET /api/v1/embed/{id}", http.StripPrefix("/api/v1/", embedHandler))

	// Post create a job to read the embedded message, if found
	extractHandler := handlers.NewExtractHandler(validator)
	mux.Handle("POST /api/v1/extract", http.StripPrefix("/api/v1/", extractHandler))

	// Get the message if one was stegg'd
	mux.Handle("GET /api/v1/extract/{id}", http.StripPrefix("/api/v1/", extractHandler))

	// Liveness check, and stats
	statsHandler := handlers.NewStatsHandler(stats)
	mux.Handle("GET /api/v1/stats", http.StripPrefix("/api/v1/", statsHandler))
	return mux
}
