package httputil

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewError(ctx *gin.Context, status int, error string) {
	er := HTTPError{
		Code:    status,
		Message: error,
	}
	ctx.JSON(status, er)
}

func NoErrorEmpty(ctx *gin.Context) {
	NoError(ctx, map[string]interface{}{})
}

func NoError(ctx *gin.Context, payload map[string]interface{}) {
	m := map[string]interface{}{
		"status": http.StatusOK,
	}
	for k, v := range payload {
		m[k] = v
	}
	ctx.JSON(http.StatusOK, m)
}

// HTTPError example
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}
