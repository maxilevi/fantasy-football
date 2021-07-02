package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/context"
	"log"
	"net/http"
	"strings"
)

var JWTKey = []byte("fantasy_football")

type AuthClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func failedAuth(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	_, err := w.Write([]byte(`{"error": true, "message": "` + message + `"}`))
	if err != nil {
		log.Printf("Failed to write response")
	}
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			failedAuth(w, http.StatusBadRequest, `{error: true, message: "Bad request"}`)
			return
		}
		tokenString := splitToken[1]
		claims := &AuthClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return JWTKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				failedAuth(w, http.StatusUnauthorized, `{error: true, message: "Unauthorized"}`)
				return
			}
			failedAuth(w, http.StatusBadRequest, `{error: true, message: "Bad request"}`)
			return
		}
		if token.Valid {
			context.Set(r, "token", token)
			context.Set(r, "email", claims.Email)
			next.ServeHTTP(w, r)
		} else {
			failedAuth(w, http.StatusUnauthorized, `{error: true, message: "Unauthorized"}`)
		}
	})
}

