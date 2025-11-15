package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/aburizalpurnama/travel/internal/app/database"
	"github.com/aburizalpurnama/travel/internal/app/domain/product"
	"github.com/aburizalpurnama/travel/internal/app/domain/user"
	"github.com/aburizalpurnama/travel/internal/app/repository"
	"github.com/aburizalpurnama/travel/internal/app/router"
	"github.com/aburizalpurnama/travel/internal/config"
	"github.com/aburizalpurnama/travel/internal/pkg/mapper"
	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db, err := database.NewGorm(cfg)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	logLevel := getLogLevel(cfg.AppEnv)
	if cfg.AppLogLevel != "" {
		logLevel = convertLogLevel(cfg.AppLogLevel)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel}))
	routerOpts := injectDependencies(db)
	routerOpts.Logger = logger

	tp, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}
	// Pastikan tracer di-flush saat aplikasi berhenti
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	app := fiber.New()
	app.Use(fiberLogger.New())

	router.SetupRoutesV1(app, routerOpts)

	port := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Fatal(app.Listen(port))
}

func injectDependencies(db *gorm.DB) *router.Option {
	userRepository := user.NewUserRepository(db)
	userService := user.NewUserService(userRepository)

	uow := repository.NewGormUnitOfWork(db)
	mapper := mapper.NewCopierMapper()

	productService := product.NewService(uow, mapper)

	userHandler := user.NewUserHandler(userService)
	productHandler := product.NewHandler(productService)

	return &router.Option{
		UserHandler:    userHandler,
		ProductHandler: productHandler,
	}
}

// initTracer menginisialisasi OpenTelemetry TracerProvider
func initTracer() (*sdktrace.TracerProvider, error) {
	// 1. Buat Exporter
	// Ganti 'stdouttrace' dengan 'jaeger', 'otlptracehttp', dll. untuk produksi
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	// 2. Buat Resource (metadata tentang aplikasi Anda)
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("travel-api-service"), // Ganti dengan nama service Anda
		attribute.String("environment", os.Getenv("APP_ENV")),
	)

	// 3. Buat Tracer Provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// 4. Set sebagai Global Provider (ini yang membuatnya "agnostik" di layer lain)
	otel.SetTracerProvider(tp)

	return tp, nil
}

func getLogLevel(v string) slog.Level {
	switch v {
	case "development":
		return slog.LevelDebug
	case "staging":
		return slog.LevelInfo
	case "production":
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}

func convertLogLevel(v string) slog.Level {
	switch v {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
