package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
)

var errReadOnlyInstance = apperrors.New(
	"READ_ONLY_INSTANCE", "this instance only accepts read requests", http.StatusServiceUnavailable,
)

func ReadOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			c.Next()
		default:
			_ = c.Error(errReadOnlyInstance.WithSource())
			c.Abort()
		}
	}
}
