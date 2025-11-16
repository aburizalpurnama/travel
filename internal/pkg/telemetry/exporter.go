package telemetry

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
)

// newExporter is an internal factory that creates and configures a trace exporter based on the provided options.
func newExporter(ctx context.Context, opt Option) (sdktrace.SpanExporter, error) {
	switch opt.Exporter {
	case "stdout":
		return stdouttrace.New(stdouttrace.WithPrettyPrint())

	case "otlp":
		// OTLP is the universal exporter standard supported by vendors like Datadog, Grafana, Sentry, etc.
		opts := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpoint(opt.OtlpEndpoint),
		}

		// Configure secure (TLS) or insecure connection
		if opt.OtlpInsecure {
			opts = append(opts, otlptracegrpc.WithInsecure())
		} else {
			opts = append(opts, otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")))
		}

		// (Optional) Add custom headers for vendor authentication or metadata
		if opt.OtlpHeaders != "" {
			headers := make(map[string]string)
			// Parse comma-separated headers (e.g., "key1=value1,key2=value2")
			for _, header := range strings.Split(opt.OtlpHeaders, ",") {
				parts := strings.SplitN(header, "=", 2)
				if len(parts) == 2 {
					headers[parts[0]] = parts[1]
				}
			}
			opts = append(opts, otlptracegrpc.WithHeaders(headers))
		}

		client := otlptracegrpc.NewClient(opts...)
		return otlptrace.New(ctx, client)

	default:
		return nil, fmt.Errorf("unknown tracer exporter: %s", opt.Exporter)
	}
}
