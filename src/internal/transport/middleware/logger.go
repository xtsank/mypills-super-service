package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		log.Printf("[INFO] %d %s %s (%v)", c.Writer.Status(), c.Request.Method, c.Request.URL.Path, time.Since(start))

		for _, err := range c.Errors {
			log.Printf("[ERROR] %v %v", err.Error(), err.Unwrap()) //TODO
		}
	}
}
