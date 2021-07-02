package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func AddUserRoutes(r *mux.Router) {
	r.HandleFunc("/users", handlePostUser).Methods( "POST")
}

func handlePostUser(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hi!!!")
}