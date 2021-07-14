package middleware

import (
	"../httputil"
	"../models"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Validates an already authenticated user has admin privileges
func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get("user")
		if !ok {
			httputil.NewError(c, http.StatusUnauthorized, "Invalid authentication")
			c.Abort()
		}
		if u, ok := v.(models.User); ok && u.IsAdmin() {
			log.Println(fmt.Sprintf("administrator role succesfully validated for request %v", c.Request.RequestURI))
			c.Next()
		} else {
			httputil.NewError(c, http.StatusUnauthorized, "User does not have administrator rights")
			c.Abort()
		}
	}
}
