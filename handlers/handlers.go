package handlers

import (
	"io"
	"net/http"
)

// Create a new job for steg, read multi-part file upload
func PostSteganographyHandler(w http.ResponseWriter, r *http.Request) {

}

// Get steg state/image
func GetSteganographyHandler(w http.ResponseWriter, r *http.Request) {

}

// Create a new job to decode
func PostDecodeHandler(w http.ResponseWriter, r *http.Request) {

}

// Get decode job/message
func GetDecodeHandler(w http.ResponseWriter, r *http.Request) {

}

// Get some generic stats
func GetStats(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, Stats")
}
