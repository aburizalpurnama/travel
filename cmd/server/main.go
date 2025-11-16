package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/aburizalpurnama/travel/internal/app/database"
	"github.com/aburizalpurnama/travel/internal/app/domain/product"
	"github.com/aburizalpurnama/travel/internal/app/repository"
	"github.com/aburizalpurnama/travel/internal/app/router"
	"github.com/aburizalpurnama/travel/internal/config"
	"github.com/aburizalpurnama/travel/internal/pkg/mapper"
	"github.com/aburizalpurnama/travel/internal/pkg/telemetry"
	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
)

func main() {
	// Load configuration from environment variables or .env file
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Configure logger level based on environment
	logLevel := getLogLevel(cfg.AppEnv)
	if cfg.AppLogLevel != "" {
		logLevel = convertLogLevel(cfg.AppLogLevel)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel}))

	// Initialize OpenTelemetry tracer provider
	shutdownTracer, err := telemetry.InitTracerProvider(telemetry.Option{
		Enabled:      cfg.Tracing.Enabled,
		ServiceName:  cfg.AppName,
		Environment:  cfg.AppEnv,
		Exporter:     cfg.Tracing.Exporter,
		OtlpEndpoint: cfg.Tracing.OtlpEndpoint,
		OtlpInsecure: cfg.Tracing.OtlpInsecure,
		OtlpHeaders:  cfg.Tracing.OtlpHeaders,
	})
	if err != nil {
		log.Fatalf("Error initializing tracer: %v", err)
	}

	// Ensure tracer provider shuts down gracefully when the application exits
	defer func() {
		if err := shutdownTracer(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// Initialize database connection
	db, err := database.NewGorm(cfg)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Inject dependencies and configure router options
	routerOpts := injectDependencies(db)
	routerOpts.Logger = logger

	// Initialize Fiber app
	app := fiber.New()
	app.Use(fiberLogger.New())

	// Setup API routes
	router.SetupRoutesV1(app, routerOpts)

	// Start the server
	port := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Fatal(app.Listen(port))
}

// injectDependencies wires up the application dependencies (repositories, services, handlers).
func injectDependencies(db *gorm.DB) *router.Option {
	uow := repository.NewGormUnitOfWork(db)
	mapper := mapper.NewCopierMapper()

	productService := product.NewService(uow, mapper)
	productHandler := product.NewHandler(productService)

	return &router.Option{
		ProductHandler: productHandler,
	}
}

// getLogLevel returns the appropriate slog.Level based on the application environment.
func getLogLevel(env string) slog.Level {
	switch env {
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

// convertLogLevel converts a string log level to slog.Level.
func convertLogLevel(level string) slog.Level {
	switch level {
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
