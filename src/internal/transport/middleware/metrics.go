package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func HTTPMetrics(registry prometheus.Registerer) gin.HandlerFunc {
	duration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "mypills_http_request_duration_seconds",
		Help:    "Duration of API requests, in seconds.",
		Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 0.75, 1, 2.5, 5},
	}, []string{"method", "route", "status"})
	registry.MustRegister(duration)

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}
		method := c.Request.Method
		switch method {
		case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodPost,
			http.MethodPut, http.MethodPatch, http.MethodDelete:
		default:
			method = "OTHER"
		}
		duration.WithLabelValues(method, route, strconv.Itoa(c.Writer.Status())).Observe(time.Since(start).Seconds())
	}
}
