package logger

import (
	"context"
	"os"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
}

// NewWithTracing creates a new logger with tracing integration.
func NewWithTracing(service string, version string) *Logger {
	config := zap.NewProductionConfig()

	// Configure log level based on environment.
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	// Add service name and structured fields for observability.
	config.InitialFields = map[string]any{
		"service": service,
		"version": version,
		"pid":     os.Getpid(),
	}

	// Configure output format for better parsing by log aggregators.
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return &Logger{logger.Sugar()}
}

// WithTrace adds tracing context to log entries.
func (l *Logger) WithTrace(ctx context.Context) *Logger {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return &Logger{
			l.SugaredLogger.With(
				"trace_id", span.SpanContext().TraceID().String(),
				"span_id", span.SpanContext().SpanID().String(),
			),
		}
	}
	return l
}

// LogRequestContext logs request details with tracing context.
func (l *Logger) LogRequestContext(
	ctx context.Context, method, path, remoteAddr string, statusCode int, duration float64,
) {
	l.WithTrace(ctx).Infow("HTTP request processed",
		"method", method,
		"path", path,
		"duration_ms", duration,
		"remote_addr", remoteAddr,
		"status_code", statusCode,
	)
}
