package router

import (
	"log/slog"

	"github.com/aburizalpurnama/travel/internal/app/domain/product"
	"github.com/aburizalpurnama/travel/internal/app/middleware"
	"github.com/gofiber/fiber/v2"
)

// Option holds the dependencies required to configure the router.
type Option struct {
	Logger         *slog.Logger
	ProductHandler *product.Handler
}

// SetupRoutesV1 configures the API routes for version 1.
func SetupRoutesV1(app *fiber.App, opt *Option) {
	// API Group v1
	api := app.Group("/api/v1")

	// Global Middleware
	api.Use(middleware.RequestLogger(opt.Logger))

	// Register domain-specific routes
	product.NewRoute(api, opt.ProductHandler)
}
