package handlers

import (
	"log"
	"net/http"
	"../repos"
)

type HandlerFunc func (w http.ResponseWriter, req *http.Request)
type HandlerFuncWithDb func (w http.ResponseWriter, req *http.Request, repo repos.Repository)

func wrap(f HandlerFuncWithDb, repo repos.Repository) HandlerFunc {
	return func (w http.ResponseWriter, req *http.Request) {
		f(w, req, repo)
	}
}

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