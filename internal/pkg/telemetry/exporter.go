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

// newExporter adalah factory internal yang memilih exporter berdasarkan config.
func newExporter(ctx context.Context, opt Option) (sdktrace.SpanExporter, error) {
	switch opt.Exporter {
	case "stdout":
		return stdouttrace.New(stdouttrace.WithPrettyPrint())

	case "otlp":
		// Ini adalah exporter universal untuk vendor seperti Datadog, Grafana, Sentry, dll.
		opts := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpoint(opt.OtlpEndpoint),
		}

		// Atur koneksi aman (TLS) atau tidak aman (insecure)
		if opt.OtlpInsecure {
			opts = append(opts, otlptracegrpc.WithInsecure())
		} else {
			opts = append(opts, otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")))
		}

		// (Opsional) Tambahkan headers kustom untuk autentikasi vendor
		if opt.OtlpHeaders != "" {
			headers := make(map[string]string)
			for header := range strings.SplitSeq(opt.OtlpHeaders, ",") {
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
