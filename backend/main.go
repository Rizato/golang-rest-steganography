package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"steg/api/v1/views"
	"steg/middleware"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	log.Printf("Listening on port 8080")
	// connect to db
	dbPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	// v1 mux
	apiV1 := views.ConfigureViews(dbPool)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/", apiV1)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Working!\n")
	})

	// add middleware
	handler := middleware.ConfigureMiddleware(mux)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
