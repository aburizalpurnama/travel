package router

import (
	"log/slog"

	"github.com/aburizalpurnama/travel/internal/app/domain/product"
	"github.com/aburizalpurnama/travel/internal/app/domain/user"
	"github.com/aburizalpurnama/travel/internal/app/middleware"
	"github.com/gofiber/fiber/v2"
)

type Option struct {
	Logger         *slog.Logger
	UserHandler    *user.UserHandler
	ProductHandler *product.Handler
}

func SetupRoutesV1(app *fiber.App, opt *Option) {
	api := app.Group("/api/v1")
	api.Use(middleware.RequestLogger(opt.Logger))

	// Rute untuk User
	user := api.Group("/users")
	user.Post("/", opt.UserHandler.CreateUser)
	user.Get("/", opt.UserHandler.GetUsers)
	user.Get("/:id", opt.UserHandler.GetUser)
	// Tambahkan rute untuk Update dan Delete

	product := api.Group("/products")
	product.Post("/", opt.ProductHandler.CreateProduct)
	product.Get("/", opt.ProductHandler.GetProducts)
	product.Get("/:id", opt.ProductHandler.GetProduct)
	product.Patch("/:id", opt.ProductHandler.UpdateProduct)
	product.Delete("/:id", opt.ProductHandler.DeleteProduct)

}
