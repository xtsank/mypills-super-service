package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResponseHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		payload, pExists := c.Get(ResponsePayloadKey)
		status, sExists := c.Get(ResponseStatusKey)

		if !pExists || !sExists {
			if c.Writer.Status() != http.StatusOK {
				return
			}
			c.Status(http.StatusNoContent)
			return
		}
		s, ok := status.(int)
		if !ok {
			s = http.StatusInternalServerError
		}
		c.JSON(s, payload)
	}
}
