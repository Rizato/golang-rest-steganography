package views

import (
	"net/http"
	"steg/api/v1/handlers"
	"steg/api/v1/models"
	"steg/api/v1/repository"
	"steg/api/v1/services"
	"steg/crud"
	"steg/middleware"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConfigureViews(dbPool *pgxpool.Pool) http.Handler {
	mux := http.NewServeMux()
	imageRepository := repository.NewImageRepository(dbPool)
	embedJobRepository := repository.NewEmbedJobRepository(dbPool)
	extractJobRepository := repository.NewExtractJobRepository(dbPool)

	imageService := services.NewImageService(imageRepository)
	embedJobService := services.NewEmbedJobService(embedJobRepository, imageRepository)
	extractJobService := services.NewExtractJobService(extractJobRepository, imageRepository)

	// Uploads and downloads are special cases due handling file content, not json
	// TODO Put strict rate limits on these routes
	mux.Handle("POST /api/v1/images", http.StripPrefix("/api/v1/", handlers.HandleUpload(imageService)))
	mux.Handle("GET /api/v1/images/{id}/download", http.StripPrefix("/api/v1/", handlers.HandleDownload(imageService)))

	// GET/DELETE files
	fileHandler := crud.NewItemHandler[models.Image](handlers.NewFileCrud(imageService))
	mux.Handle("/api/v1/images/{id}", http.StripPrefix("/api/v1/", fileHandler))

	embedCrud := handlers.NewEmbedJobCrud(embedJobService)
	// POST/GET to create or list
	embedListHandler := crud.NewListHandler[models.EmbedJob](embedCrud)
	mux.Handle("/api/v1/embed", http.StripPrefix("/api/v1/", embedListHandler))

	// GET/DELETE to access or delete an item
	embedJobHandler := crud.NewItemHandler[models.EmbedJob](embedCrud)
	mux.Handle("/api/v1/embed/{id}", http.StripPrefix("/api/v1/", embedJobHandler))

	// POST to trigger the processing
	embedStartHandler := handlers.NewJobStartHandler(embedJobService)
	mux.Handle("POST /api/v1/embed/{id}/start", http.StripPrefix("/api/v1/", embedStartHandler))

	extractCrud := handlers.NewExtractJobCrud(extractJobService)
	// POST/GET to create or list
	extractListHandler := crud.NewListHandler[models.ExtractJob](extractCrud)
	mux.Handle("/api/v1/extract", http.StripPrefix("/api/v1/", extractListHandler))

	// GET/DELETE to access or delete an item
	extractJobHandler := crud.NewItemHandler[models.ExtractJob](extractCrud)
	mux.Handle("/api/v1/extract/{id}", http.StripPrefix("/api/v1/", extractJobHandler))

	// POST to trigger the processing
	extractStartHandler := handlers.NewJobStartHandler(extractJobService)
	mux.Handle("POST /api/v1/extract/{id}/start", http.StripPrefix("/api/v1/", extractStartHandler))

	// V1 stats
	stats := middleware.NewStats()
	statsService := services.NewStatsService(imageRepository, embedJobRepository, extractJobRepository, stats)
	statsHandler := handlers.NewStatsHandler(statsService)
	mux.Handle("GET /api/v1/stats", http.StripPrefix("/api/v1/", statsHandler))

	// Middleware to collect v1 stats
	return middleware.NewStatsMiddleware(stats).Wrap(mux)
}
