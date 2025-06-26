package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	DatabaseConnections   prometheus.Gauge
	ActiveRequests        prometheus.Gauge
	HTTPRequestsTotal     *prometheus.CounterVec
	HTTPRequestDuration   *prometheus.HistogramVec
	DatabaseQueryDuration *prometheus.HistogramVec
}

func New() *Metrics {
	return &Metrics{
		// HTTP request metrics
		HTTPRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status_code"},
		),
		HTTPRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),

		// Database metrics
		DatabaseConnections: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "database_connections_active",
				Help: "Number of active database connections",
			},
		),
		DatabaseQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "database_query_duration_seconds",
				Help:    "Duration of database queries in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"query_type"},
		),

		// Application metrics
		ActiveRequests: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "active_requests",
				Help: "Number of requests currently being processed",
			},
		),
	}
}

// RecordHTTPRequest records metrics for an HTTP request.
func (m *Metrics) RecordHTTPRequest(method, endpoint, statusCode string, duration float64) {
	m.HTTPRequestsTotal.WithLabelValues(method, endpoint, statusCode).Inc()
	m.HTTPRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
}

// RecordDatabaseQuery records metrics for a database query.
func (m *Metrics) RecordDatabaseQuery(queryType string, duration float64) {
	m.DatabaseQueryDuration.WithLabelValues(queryType).Observe(duration)
}
