package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Counter untuk menghitung jumlah permintaan
var requestCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	},
	[]string{"method", "endpoint"},
)

// Counter untuk menghitung jumlah kesalahan
var errorCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_errors_total",
		Help: "Total number of HTTP errors",
	},
	[]string{"method", "endpoint"},
)

// Histogram untuk mengukur latensi permintaan
var requestLatency = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_latency_seconds",
		Help:    "Request latency in seconds",
		Buckets: prometheus.LinearBuckets(0.1, 0.1, 10), // Rentang bucket 0.1 detik hingga 1 detik
	},
	[]string{"method"},
)

// Gauge untuk menghitung jumlah koneksi aktif
var activeConnections = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "active_connections",
		Help: "Number of active connections",
	},
)

func main() {
	// Inisialisasi Prometheus
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(errorCount)
	prometheus.MustRegister(requestLatency)
	prometheus.MustRegister(activeConnections)

	// Membuat router Gin
	router := gin.Default()

	// Endpoint untuk menambahkan tugas baru
	router.POST("/todo", func(c *gin.Context) {
		// Logika untuk menambahkan tugas
		// ...

		endpoint := c.FullPath() // Get the endpoint path
		method := "GET"          // Set the HTTP method

		// Increment request metric in Prometheus with method and endpoint labels
		requestCount.WithLabelValues(method, endpoint).Inc()

		c.JSON(200, gin.H{"message": "Task added successfully"})
	})

	// Endpoint untuk mengambil semua tugas
	router.GET("/todo", func(c *gin.Context) {
		// Logika untuk mengambil tugas
		// ...

		endpoint := c.FullPath() // Get the endpoint path
		method := "GET"          // Set the HTTP method

		// Increment request metric in Prometheus with method and endpoint labels
		requestCount.WithLabelValues(method, endpoint).Inc()

		// Menambahkan metrik ke Prometheus
		// requestCount.WithLabelValues("GET").Inc()

		c.JSON(200, gin.H{"tasks": []string{"Task 1", "Task 2"}})
	})

	// Endpoint to delete a task
	router.DELETE("/todo/:id", func(c *gin.Context) {
		// Logic to delete a task
		// ...

		endpoint := c.FullPath() // Get the endpoint path
		method := "DELETE"       // Set the HTTP method

		// Increment request metric in Prometheus with method and endpoint labels
		requestCount.WithLabelValues(method, endpoint).Inc()

		c.JSON(200, gin.H{"message": "Task deleted successfully"})
	})

	// Middleware untuk menangkap kesalahan
	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := float64(time.Since(start).Seconds())

		// Menambahkan metrik latensi ke Prometheus
		requestLatency.WithLabelValues(c.Request.Method).Observe(latency)
		activeConnections.Set(float64(getActiveConnections()))

		// Menambahkan metrik kesalahan ke Prometheus
		if c.Writer.Status() >= 400 {
			errorCount.WithLabelValues(c.Request.Method, c.Request.URL.Path).Inc()
		}
	})

	// Endpoint Prometheus untuk mengumpulkan metrik
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Menjalankan server
	router.Run(":8081")
}

// Fungsi untuk mendapatkan jumlah koneksi aktif (contoh fiktif)
func getActiveConnections() int {
	// Logika untuk menghitung jumlah koneksi aktif
	// ...

	return 10 // Contoh: Mengembalikan jumlah koneksi aktif (misalnya 10)
}
