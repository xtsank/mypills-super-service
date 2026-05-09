package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		var appErr *apperrors.AppError
		if !errors.As(err, &appErr) {
			appErr = apperrors.ErrInternal.WithError(err)
		}

		c.Set(ResponsePayloadKey, appErr)
		c.Set(ResponseStatusKey, appErr.HTTPStatus)
	}
}
