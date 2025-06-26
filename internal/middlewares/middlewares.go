package middlewares

import (
	"net/http"
)

// responseWriter wraps http.ResponseWriter to capture status code.
type responseWriter struct {
	statusCode int
	http.ResponseWriter
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
