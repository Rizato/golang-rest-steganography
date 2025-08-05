package main

import (
	"log"
	"net/http"
	"steg/api/v1"
)

func main() {
	server := http.NewServeMux()
	log.Printf("Listening on port 808")
	v1.ConfigureServer(server)
	log.Fatal(http.ListenAndServe(":8080", server))
}
