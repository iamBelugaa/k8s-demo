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

func New(config *TracingConfig) (func(context.Context) error, error) {
	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(config.JaegerEndpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(config.ServiceName),
		semconv.ServiceVersionKey.String(config.ServiceVersion),
		semconv.DeploymentEnvironmentKey.String(config.Environment),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(
			exporter,
			sdktrace.WithMaxExportBatchSize(512),
			sdktrace.WithBatchTimeout(time.Second*5),
			sdktrace.WithExportTimeout(time.Second*30),
		),
		sdktrace.WithResource(resource),
		sdktrace.WithSampler(getSamplerForEnvironment(config.Environment)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)

	return tp.Shutdown, nil
}

func GetTracer(name string) trace.Tracer {
	return otel.Tracer(name)
}

func StartSpan(
	ctx context.Context, tracerName, spanName string, options ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	tracer := GetTracer(tracerName)
	return tracer.Start(ctx, spanName, options...)
}

func getSamplerForEnvironment(environment string) sdktrace.Sampler {
	switch environment {
	case "PRODUCTION":
		return sdktrace.TraceIDRatioBased(0.2)
	case "DEVELOPMENT":
		return sdktrace.AlwaysSample()
	default:
		return sdktrace.TraceIDRatioBased(0.5)
	}
}
