package middleware

import (
	"../models"
	"fmt"
	"github.com/gin-gonic/gin"
	"../httputil"
	"log"
	"net/http"
)

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
