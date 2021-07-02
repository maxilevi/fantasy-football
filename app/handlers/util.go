package handlers

import (
	"log"
	"net/http"
)

func writeResponse(w http.ResponseWriter, code int, message [] byte) {
	w.WriteHeader(code)
	_, err := w.Write(message)
	if err != nil {
		log.Printf("Failed to write response")
	}
}

func writeError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	_, err := w.Write([]byte(`{"error": true, "message": "` + message + `"}`))
	if err != nil {
		log.Printf("Failed to write response")
	}
}