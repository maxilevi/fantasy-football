package controllers

import (
	"../models"
	"fmt"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func noError() []byte {
	return []byte(`{"error": false}`)
}

func writeResponse(w http.ResponseWriter, code int, message []byte) {
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

func getAuthenticatedUserFromRequest(w http.ResponseWriter, req *http.Request) (models.User, error) {
	val, ok := context.GetOk(req, "user")
	if !ok {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return models.User{}, fmt.Errorf("internal server error")
	}
	user, ok := val.(models.User)
	if !ok {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return models.User{}, fmt.Errorf("internal server error")
	}
	return user, nil
}

func parseIdFromRequest(w http.ResponseWriter, req *http.Request) (uint, error) {
	vars := mux.Vars(req)
	id, ok := strconv.ParseInt(vars["id"], 10, 32)
	if ok != nil {
		writeError(w, http.StatusBadRequest, "Invalid id")
		return 0, fmt.Errorf("invalid id")
	}
	return uint(id), nil
}