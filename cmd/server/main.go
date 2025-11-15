package main

import (
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
	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
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
	productService := product.NewService(uow)

	userHandler := user.NewUserHandler(userService)
	productHandler := product.NewHandler(productService)

	return &router.Option{
		UserHandler:    userHandler,
		ProductHandler: productHandler,
	}
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
