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

type Option struct {
	Enabled      bool
	ServiceName  string
	Environment  string
	Exporter     string
	OtlpEndpoint string
	OtlpInsecure bool
	OtlpHeaders  string
}

// InitTracerProvider menginisialisasi provider OTel dan mendaftarkannya secara global.
// Ia mengembalikan fungsi shutdown yang harus dipanggil saat aplikasi berhenti.
func InitTracerProvider(opt Option) (func(context.Context) error, error) {
	// 1. Cek apakah tracing diaktifkan
	if !opt.Enabled {
		slog.Info("Tracing is DISABLED (No-op provider).")
		// Kembalikan fungsi shutdown "no-op" (kosong)
		return func(context.Context) error { return nil }, nil
	}

	slog.Info("Tracing is ENABLED", "exporter", opt.Exporter)

	// 2. Buat exporter yang dikonfigurasi menggunakan factory
	exporter, err := newExporter(context.Background(), opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	// 3. Buat Resource (metadata tentang aplikasi Anda)
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(opt.ServiceName),
		attribute.String("environment", opt.Environment),
	)

	// 4. Buat Tracer Provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// 5. Set sebagai Global Provider
	otel.SetTracerProvider(tp)

	// 6. Kembalikan fungsi shutdown yang sebenarnya
	return tp.Shutdown, nil
}
