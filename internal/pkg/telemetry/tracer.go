package telemetry

import (
	"context"
	"fmt"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// Option holds the configuration parameters for initializing OpenTelemetry tracing.
type Option struct {
	Enabled      bool
	ServiceName  string
	Environment  string
	Exporter     string
	OtlpEndpoint string
	OtlpInsecure bool
	OtlpHeaders  string
}

// InitTracerProvider initializes the OpenTelemetry tracer provider and registers it globally.
// It returns a shutdown function that must be called when the application terminates to flush pending spans.
func InitTracerProvider(opt Option) (func(context.Context) error, error) {
	if !opt.Enabled {
		slog.Info("Tracing is DISABLED (No-op provider).")
		// Return a no-op shutdown function
		return func(context.Context) error { return nil }, nil
	}

	slog.Info("Tracing is ENABLED", "exporter", opt.Exporter)

	// Create an exporter configured via the internal factory
	exporter, err := newExporter(context.Background(), opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	// Create a Resource (metadata about your application)
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(opt.ServiceName),
		attribute.String("environment", opt.Environment),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}
