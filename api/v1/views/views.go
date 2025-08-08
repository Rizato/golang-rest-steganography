package views

import (
	"net/http"
	"steg/api/v1/handlers"
	"steg/api/v1/models"
	"steg/middleware"
)

func ConfigureViews(stats *middleware.Stats) http.Handler {
	mux := http.NewServeMux()
	ds := models.NewDatastore()

	// Uploads and downloads are special cases due handling file content, not json
	// TODO Put strict rate limits on these routes
	mux.Handle("POST /api/v1/images", http.StripPrefix("/api/v1/", handlers.HandleUpload(ds)))
	mux.Handle("GET /api/v1/images/{id}", http.StripPrefix("/api/v1/", handlers.HandleDownload(ds)))

	// File crud, get the metadata, and support deletion
	fileHandler := handlers.NewFileCrudHandler(ds)
	mux.Handle("GET /api/v1/images/{id}/metadata", http.StripPrefix("/api/v1/", fileHandler))
	mux.Handle("DELETE /api/v1/images/{id}", http.StripPrefix("/api/v1/", fileHandler))

	embedHandler := handlers.NewEmbedCrudHandler(ds)
	// POST create a new job to embed a message in the given image
	mux.Handle("POST /api/v1/embed", http.StripPrefix("/api/v1/", embedHandler))

	// Get details about a job, and the url of the embedded image if finished
	mux.Handle("GET /api/v1/embed/{id}", http.StripPrefix("/api/v1/", embedHandler))

	// Trigger the processing
	//mux.Handle("POST /api/v1/embed/{id}/start", http.StripPrefix("/api/v1/", ))

	extractHandler := handlers.NewExtractCrudHandler(ds)
	// Post create a job to read the embedded message, if found
	mux.Handle("POST /api/v1/extract", http.StripPrefix("/api/v1/", extractHandler))

	// Get the message if the image had an embedded message, and it is finished
	mux.Handle("GET /api/v1/extract/{id}", http.StripPrefix("/api/v1/", extractHandler))

	// Trigger the processing
	//mux.Handle("POST /api/v1/extract/{id}/start", http.StripPrefix("/api/v1/", ))

	// Liveness check, and stats
	statsHandler := handlers.NewStatsHandler(stats)
	mux.Handle("GET /api/v1/stats", http.StripPrefix("/api/v1/", statsHandler))
	return mux
}
