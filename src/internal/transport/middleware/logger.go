package middleware

import (
	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		log.Printf("[INFO] %d %s %s (%v)", c.Writer.Status(), c.Request.Method, c.Request.URL.Path, time.Since(start))

		if len(c.Errors) > 0 {
			lastErr := c.Errors.Last()

			underlyingErr := errors.Unwrap(lastErr.Err)

			if underlyingErr != nil {
				log.Printf("[ERROR] Business error: '%v', Underlying system error: '%v'", lastErr.Err, underlyingErr)
			} else {
				log.Printf("[ERROR] Business error: '%v'", lastErr.Err)
			}
		}
	}
}
