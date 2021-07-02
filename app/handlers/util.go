package handlers

import (
	"gorm.io/gorm"
	"log"
	"net/http"
)

type HandlerFunc func (w http.ResponseWriter, req *http.Request)
type HandlerFuncWithDb func (w http.ResponseWriter, req *http.Request, db *gorm.DB)

func wrap(f HandlerFuncWithDb, db *gorm.DB) HandlerFunc {
	return func (w http.ResponseWriter, req *http.Request) {
		f(w, req, db)
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