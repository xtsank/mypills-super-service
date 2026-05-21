package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/xtsank/mypills-super-service/src/internal/errors"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, ok := c.Get(IsAdminKey)
		if !ok {
			_ = c.Error(errors.ErrUnauthorized.WithSource())
			c.Abort()
			return
		}

		allowed, ok := isAdmin.(bool)
		if !ok || !allowed {
			_ = c.Error(errors.ErrUnauthorized.WithSource())
			c.Abort()
			return
		}

		c.Next()
	}
}

