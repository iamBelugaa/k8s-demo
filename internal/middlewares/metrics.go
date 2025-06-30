package middlewares

import (
	"net/http"
	"strconv"
	"time"

	"github.com/iamBelugaa/k8s-demo/internal/metrics"
)

func MetricsMiddleware(metrics *metrics.Metrics) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			metrics.ActiveRequests.Inc()
			defer metrics.ActiveRequests.Dec()

			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(wrapped, r)

			statusCode := strconv.Itoa(wrapped.statusCode)
			duration := float64(time.Since(start).Milliseconds())
			metrics.RecordHTTPRequest(r.Method, r.URL.Path, statusCode, duration)
		})
	}
}
