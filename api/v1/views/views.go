package views

import (
	"net/http"
	"steg/api/v1/handlers"
	"steg/api/v1/models"
	"steg/crud"
	"steg/middleware"
)

func ConfigureViews(stats *middleware.Stats) http.Handler {
	mux := http.NewServeMux()
	ds := models.NewDatastore()

	// Uploads and downloads are special cases due handling file content, not json
	// TODO Put strict rate limits on these routes
	mux.Handle("POST /api/v1/images", http.StripPrefix("/api/v1/", handlers.HandleUpload(ds)))
	mux.Handle("GET /api/v1/images/{id}/download", http.StripPrefix("/api/v1/", handlers.HandleDownload(ds)))

	// GET/DELETE files
	fileHandler := crud.NewItemHandler[*models.ServerFile](handlers.NewFileCrud(ds))
	mux.Handle("/api/v1/images/{id}", http.StripPrefix("/api/v1/", fileHandler))

	embedCrud := handlers.NewEmbedJobCrud(ds)
	// POST/GET to create or list
	embedListHandler := crud.NewListHandler[*models.EmbedJob](embedCrud)
	mux.Handle("/api/v1/embed", http.StripPrefix("/api/v1/", embedListHandler))

	// GET/DELETE to access or delete an item
	embedJobHandler := crud.NewItemHandler[*models.EmbedJob](embedCrud)
	mux.Handle("/api/v1/embed/{id}", http.StripPrefix("/api/v1/", embedJobHandler))

	// POST to trigger the processing
	embedStartHandler := handlers.NewJobStartHandler[*models.EmbedJob](embedCrud)
	mux.Handle("POST /api/v1/embed/{id}/start", http.StripPrefix("/api/v1/", embedStartHandler))

	extractCrud := handlers.NewExtractJobCrud(ds)
	// POST/GET to create or list
	extractListHandler := crud.NewListHandler[*models.ExtractJob](extractCrud)
	mux.Handle("/api/v1/extract", http.StripPrefix("/api/v1/", extractListHandler))

	// GET/DELETE to access or delete an item
	extractJobHandler := crud.NewItemHandler[*models.ExtractJob](extractCrud)
	mux.Handle("/api/v1/extract/{id}", http.StripPrefix("/api/v1/", extractJobHandler))

	// POST to trigger the processing
	extractStartHandler := handlers.NewJobStartHandler[*models.ExtractJob](extractCrud)
	mux.Handle("POST /api/v1/extract/{id}/start", http.StripPrefix("/api/v1/", extractStartHandler))

	// Liveness check, and stats
	statsHandler := handlers.NewStatsHandler(stats, ds)
	mux.Handle("GET /api/v1/stats", http.StripPrefix("/api/v1/", statsHandler))
	return mux
}
