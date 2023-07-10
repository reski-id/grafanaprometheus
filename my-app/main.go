package main

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Counter to count the number of requests
var requestCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	},
	[]string{"method", "endpoint"},
)

// Counter to count the number of errors
var errorCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_errors_total",
		Help: "Total number of HTTP errors",
	},
	[]string{"method", "endpoint"},
)

// Histogram to measure request latency
var requestLatency = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_latency_seconds",
		Help:    "Request latency in seconds",
		Buckets: prometheus.LinearBuckets(0.1, 0.1, 10),
	},
	[]string{"method", "endpoint"},
)

// Gauge to track memory usage
var memoryUsage = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "app_memory_usage",
		Help: "Memory usage of the application",
	},
	[]string{"endpoint"},
)

// Gauge to count active connections
var activeConnections = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "active_connections",
		Help: "Number of active connections",
	},
)

// Counter for tracking the total number of scrapes by HTTP status code
var scrapeRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "promhttp_metric_handler_requests_total",
		Help: "Total number of scrapes by HTTP status code.",
	},
	[]string{"code"},
)

func main() {
	// Initialize Prometheus
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(errorCount)
	prometheus.MustRegister(requestLatency)
	prometheus.MustRegister(activeConnections)
	prometheus.MustRegister(memoryUsage)
	prometheus.MustRegister(scrapeRequestsTotal)

	// Create a Gin router
	router := gin.Default()

	// Add middleware to capture errors and collect metrics
	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := float64(time.Since(start).Seconds())
		code := strconv.Itoa(c.Writer.Status())

		method := c.Request.Method
		endpoint := c.FullPath()

		// Increment request count metric with method and endpoint labels
		requestCount.WithLabelValues(method, endpoint).Inc()

		// Add latency metric to Prometheus with method and endpoint labels
		requestLatency.WithLabelValues(method, endpoint).Observe(latency)

		// Add memory usage metric to Prometheus with endpoint label
		memoryUsage.WithLabelValues(endpoint).Set(float64(getMemoryUsage()))

		activeConnections.Set(float64(getActiveConnections()))

		// Increment the scrape requests metric with the corresponding status code
		scrapeRequestsTotal.WithLabelValues(code).Inc()

		// Increment error count metric in Prometheus with method and endpoint labels
		if c.Writer.Status() >= 400 {
			errorCount.WithLabelValues(method, endpoint).Inc()
		}
	})

	// Add routes
	router.GET("/todo", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "GET /todo"})
	})

	router.POST("/todo", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "POST /todo"})
	})

	router.DELETE("/todo/:id", func(c *gin.Context) {
		start := time.Now()
		id := c.Param("id") // Get the actual ID value from the request

		// Logic to delete the task with the given ID
		// ...

		endpoint := "/todo/" + id // Set the endpoint path with the real ID
		method := "DELETE"        // Set the HTTP method

		// Increment request metric in Prometheus with method and endpoint labels
		requestCount.WithLabelValues(method, endpoint).Inc()

		// Add latency metric to Prometheus with method and endpoint labels
		requestLatency.WithLabelValues(method, endpoint).Observe(float64(time.Since(start).Seconds()))

		// Add memory usage metric to Prometheus with endpoint label
		memoryUsage.WithLabelValues(endpoint).Set(float64(getMemoryUsage()))

		activeConnections.Set(float64(getActiveConnections()))

		// Increment error count metric in Prometheus with method and endpoint labels
		if c.Writer.Status() >= 400 {
			errorCount.WithLabelValues(method, endpoint).Inc()
		}

		c.JSON(200, gin.H{"message": "Task deleted successfully", "id": id})
	})

	// Add endpoint for Prometheus to scrape metrics
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Fallback route to handle undefined endpoints
	router.NoRoute(func(c *gin.Context) {
		endpoint := c.Request.URL.Path // Get the endpoint path
		method := c.Request.Method     // Get the HTTP method

		// Increment request metric in Prometheus with method and endpoint labels
		requestCount.WithLabelValues(method, endpoint).Inc()
		errorCount.WithLabelValues(method, endpoint).Inc()

		c.JSON(404, gin.H{"message": "Endpoint not found"})
	})

	// Run the server
	router.Run(":8081")

}

func getActiveConnections() int {
	// Logic to count active connections
	return 10
}

func getMemoryUsage() int {
	// Logic to calculate memory usage
	return 100
}
