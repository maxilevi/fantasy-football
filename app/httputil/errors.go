package httputil

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Write an error to the response using the status code and an error message
func NewError(ctx *gin.Context, status int, error string) {
	er := HTTPError{
		Code:    status,
		Message: error,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(status, er)
}

// Write an OK response with no message
func NoErrorEmpty(ctx *gin.Context) {
	NoError(ctx, map[string]interface{}{})
}

// Write an OK response with a given message
func NoError(ctx *gin.Context, payload interface{}) {
	m := map[string]interface{}{
		"code": http.StatusOK,
	}
	body, _ := json.Marshal(payload)
	payload = json.Unmarshal(body, &m)
	for k, v := range m {
		m[k] = v
	}
	ctx.JSON(http.StatusOK, m)
}

// HTTPError example
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
} // @name HTTPError
