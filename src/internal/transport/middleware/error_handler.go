package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	app_errors "github.com/xtsank/mypills-super-service/internal/errors"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		var appErr *app_errors.AppError
		if !errors.As(err, &appErr) {
			appErr = app_errors.ErrInternal.WithError(err)
		}

		c.Set(ResponsePayloadKey, appErr)
		c.Set(ResponseStatusKey, appErr.HTTPStatus)
	}
}
