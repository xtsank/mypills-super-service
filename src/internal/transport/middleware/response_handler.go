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
			c.Status(http.StatusNoContent)
			return
		}

		c.JSON(status.(int), payload)
	}
}
