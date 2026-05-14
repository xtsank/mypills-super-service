package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/service"
)

func TokenVerifier(i do.Injector) gin.HandlerFunc {
	tokenManager := do.MustInvoke[service.TokenManager](i)

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			_ = c.Error(errors.ErrUnauthorized.WithSource())
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			_ = c.Error(errors.ErrUnauthorized.WithSource())
			c.Abort()
			return
		}

		userID, isAdmin, err := tokenManager.VerifyToken(parts[1])
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		c.Set(UserIDKey, userID)
		c.Set(IsAdminKey, isAdmin)

		c.Next()
	}
}
