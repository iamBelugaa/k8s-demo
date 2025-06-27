package tracing

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

type TracingConfig struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	JaegerEndpoint string
}

// New initializes OpenTelemetry tracing with OTLP exporter.
func New(config *TracingConfig) (func(context.Context) error, error) {
	// Create OTLP HTTP exporter.
	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(config.JaegerEndpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	// Create resource with service information.
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(config.ServiceName),
		semconv.ServiceVersionKey.String(config.ServiceVersion),
		semconv.DeploymentEnvironmentKey.String(config.Environment),
	)

	// Create trace provider with batch span processor.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(
			exporter,
			sdktrace.WithMaxExportBatchSize(512),       // Maximum spans per batch.
			sdktrace.WithBatchTimeout(time.Second*5),   // Send batch every 5 seconds.
			sdktrace.WithExportTimeout(time.Second*30), // Timeout for export operations.
		),
		sdktrace.WithResource(resource),
		sdktrace.WithSampler(getSamplerForEnvironment(config.Environment)),
	)

	// Set the global trace provider.
	otel.SetTracerProvider(tp)

	// Set up propagation for distributed tracing.
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)

	return tp.Shutdown, nil
}

// GetTracer returns a tracer for the specified component.
func GetTracer(name string) trace.Tracer {
	return otel.Tracer(name)
}

// StartSpan is a helper function to start a span with common attributes.
func StartSpan(
	ctx context.Context, tracerName, spanName string, options ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	tracer := GetTracer(tracerName)
	return tracer.Start(ctx, spanName, options...)
}

// getSamplerForEnvironment returns appropriate sampling strategy based on environment.
func getSamplerForEnvironment(environment string) sdktrace.Sampler {
	switch environment {
	case "PRODUCTION":
		// In production, sample 20% of traces to balance observability with performance.
		return sdktrace.TraceIDRatioBased(0.2)
	case "DEVELOPMENT":
		// In development, sample 100% for complete debugging visibility.
		return sdktrace.AlwaysSample()
	default:
		// In other environment, sample 50%.
		return sdktrace.TraceIDRatioBased(0.5)
	}
}
