package middleware

import (
	"../httputil"
	"../repos"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"os"
	"strings"
)

// Get the secret from the environmental variables
var JWTKey = []byte(os.Getenv("JWT_SECRET"))

type AuthClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// Validates the request is authenticated with a valid user
func Auth(repo repos.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(c.GetHeader("Authorization")) == 0 {
			httputil.NewError(c, http.StatusUnauthorized, "Authorization is a required header")
			c.Abort()
			return
		}

		// Verify JWT token
		authHeader := c.GetHeader("Authorization")
		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			httputil.NewError(c, http.StatusUnauthorized, "Authorization is a required header")
			c.Abort()
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
			httputil.NewError(c, http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		if token.Valid {
			c.Set("token", token)
			c.Set("email", claims.Email)
			user, err := repo.GetUserByEmail(claims.Email)
			if err != nil {
				httputil.NewError(c, http.StatusUnauthorized, "Invalid token")
				c.Abort()
			}
			// Save the user so the handlers can use it
			c.Set("user", user)
			log.Println(fmt.Sprintf("user %v succesfully authenticated for request %v", claims.Email, c.Request.RequestURI))
			c.Next()
		} else {
			httputil.NewError(c, http.StatusUnauthorized, "Invalid token")
			c.Abort()
		}
	}
}
