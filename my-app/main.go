package main

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "http_requests_total", Help: "Total number of HTTP requests"},
		[]string{"method", "endpoint"},
	)

	errorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "http_errors_total", Help: "Total number of HTTP errors"},
		[]string{"method", "endpoint"},
	)

	requestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{Name: "http_request_latency_seconds", Help: "Request latency in seconds", Buckets: prometheus.LinearBuckets(0.1, 0.1, 10)},
		[]string{"method", "endpoint"},
	)

	activeConnections = prometheus.NewGauge(prometheus.GaugeOpts{Name: "active_connections", Help: "Number of active connections"})

	scrapeRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "promhttp_metric_handler_requests_total", Help: "Total number of scrapes by HTTP status code."},
		[]string{"code"},
	)
)

func main() {
	prometheus.MustRegister(requestCount, errorCount, requestLatency, activeConnections, scrapeRequestsTotal)

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := float64(time.Since(start).Seconds())
		code := strconv.Itoa(c.Writer.Status())

		method := c.Request.Method
		endpoint := c.FullPath()

		requestCount.WithLabelValues(method, endpoint).Inc()
		requestLatency.WithLabelValues(method, endpoint).Observe(latency)
		activeConnections.Set(float64(getActiveConnections()))
		scrapeRequestsTotal.WithLabelValues(code).Inc()

		if c.Writer.Status() >= 400 {
			errorCount.WithLabelValues(method, endpoint).Inc()
		}
	})

	router.GET("/todo", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "GET /todo"})
	})

	router.POST("/todo", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "POST /todo"})
	})

	router.DELETE("/todo/:id", func(c *gin.Context) {
		start := time.Now()
		id := c.Param("id")
		endpoint := c.Request.URL.Path
		method := c.Request.Method

		requestCount.WithLabelValues(method, endpoint).Inc()
		requestLatency.WithLabelValues(method, endpoint).Observe(float64(time.Since(start).Seconds()))
		activeConnections.Set(float64(getActiveConnections()))

		if c.Writer.Status() >= 400 {
			errorCount.WithLabelValues(method, endpoint).Inc()
		}

		c.JSON(200, gin.H{"message": "Task deleted successfully", "id": id})
	})

	router.NoRoute(func(c *gin.Context) {
		endpoint := c.Request.URL.Path
		method := c.Request.Method

		requestCount.WithLabelValues(method, endpoint).Inc()
		errorCount.WithLabelValues(method, endpoint).Inc()

		c.JSON(404, gin.H{"message": "Endpoint not found"})
	})

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.Run(":8081")
}

func getActiveConnections() int {
	return 1
}
