package handlers

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

func AddSessionRoutes(r *mux.Router, db *gorm.DB) {
	r.HandleFunc("/session", wrap(handlePostSession, db)).Methods( "POST")
}

func handlePostSession(w http.ResponseWriter, req *http.Request, db *gorm.DB) {

}