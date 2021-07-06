package middleware

import (
	"../models"
	"fmt"
	"github.com/gorilla/context"
	"log"
	"net/http"
)

func Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v, ok := context.GetOk(r, "user")
		if !ok {
			failedAuth(w, http.StatusUnauthorized, "Unauthorized")
		}
		if u, ok := v.(models.User); ok && u.IsAdmin() {
			log.Println(fmt.Sprintf("administrator role succesfully validated for request %v", r.RequestURI))
			next.ServeHTTP(w, r)
		} else {
			failedAuth(w, http.StatusUnauthorized, "Unauthorized")
		}
	})
}
