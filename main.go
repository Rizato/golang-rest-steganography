package main

import (
	"io"
	"log"
	"net/http"
	"steg/api/v1/views"
	"steg/middleware"
)

func main() {

	log.Printf("Listening on port 8080")
	stats := middleware.NewStats()
	// v1 mux
	apiV1 := views.ConfigureViews(stats)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/", apiV1)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Working!\n")
	})

	// add middleware
	handler := middleware.ConfigureMiddleware(mux, stats)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
