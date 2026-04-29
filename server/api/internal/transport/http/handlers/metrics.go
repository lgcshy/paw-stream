package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// Metrics holds Prometheus metrics for the API
type Metrics struct {
	httpRequestsTotal   *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
	authFailuresTotal   prometheus.Counter
	activeStreams        prometheus.Gauge
	devicesTotal        *prometheus.GaugeVec
}

// NewMetrics creates and registers Prometheus metrics
func NewMetrics() *Metrics {
	m := &Metrics{
		httpRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "pawstream_http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		httpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "pawstream_http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
		authFailuresTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "pawstream_auth_failures_total",
				Help: "Total number of authentication failures",
			},
		),
		activeStreams: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "pawstream_active_streams",
				Help: "Number of currently active streams",
			},
		),
		devicesTotal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "pawstream_devices_total",
				Help: "Total number of devices by status",
			},
			[]string{"status"},
		),
	}

	prometheus.MustRegister(
		m.httpRequestsTotal,
		m.httpRequestDuration,
		m.authFailuresTotal,
		m.activeStreams,
		m.devicesTotal,
	)

	return m
}

// Middleware returns a Fiber middleware that records request metrics
func (m *Metrics) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Response().StatusCode())
		path := c.Route().Path
		method := c.Method()

		m.httpRequestsTotal.WithLabelValues(method, path, status).Inc()
		m.httpRequestDuration.WithLabelValues(method, path).Observe(duration)

		// Track auth failures
		if c.Response().StatusCode() == 401 || c.Response().StatusCode() == 403 {
			m.authFailuresTotal.Inc()
		}

		return err
	}
}

// IncrementActiveStreams increments the active stream counter
func (m *Metrics) IncrementActiveStreams() {
	m.activeStreams.Inc()
}

// DecrementActiveStreams decrements the active stream counter
func (m *Metrics) DecrementActiveStreams() {
	m.activeStreams.Dec()
}

// SetDeviceCounts sets the device count gauges
func (m *Metrics) SetDeviceCounts(online, offline int) {
	m.devicesTotal.WithLabelValues("online").Set(float64(online))
	m.devicesTotal.WithLabelValues("offline").Set(float64(offline))
}

// Handler returns a Fiber handler for the /metrics endpoint
func MetricsHandler() fiber.Handler {
	handler := promhttp.Handler()
	return func(c *fiber.Ctx) error {
		fasthttpadaptor.NewFastHTTPHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.ServeHTTP(w, r)
		}))(c.Context())
		return nil
	}
}
