package main

import (
	"log"
	"net/http"
	"steg/api/v1/views"
)

func main() {
	server := http.NewServeMux()
	log.Printf("Listening on port 808")
	views.ConfigureViews(server)
	log.Fatal(http.ListenAndServe(":8080", server))
}
