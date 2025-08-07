package main

import (
	"log"
	"net/http"
	"steg/views/api/v1"
)

func main() {
	server := http.NewServeMux()
	log.Printf("Listening on port 808")
	v1.ConfigureViews(server)
	log.Fatal(http.ListenAndServe(":8080", server))
}
