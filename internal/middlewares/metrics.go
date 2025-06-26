package middlewares

import (
	"net/http"
	"strconv"
	"time"

	"github.com/iamNilotpal/k8s-demo/internal/metrics"
)

// MetricsMiddleware records HTTP metrics for each request.
func MetricsMiddleware(metrics *metrics.Metrics) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Increment active requests.
			metrics.ActiveRequests.Inc()
			defer metrics.ActiveRequests.Dec()

			// Wrap response writer to capture status code.
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Process request.
			next.ServeHTTP(wrapped, r)

			// Record metrics.
			statusCode := strconv.Itoa(wrapped.statusCode)
			duration := float64(time.Since(start).Milliseconds())

			metrics.RecordHTTPRequest(r.Method, r.URL.Path, statusCode, duration)
		})
	}
}
