package router

import (
	"log/slog"

	"github.com/aburizalpurnama/travel/internal/app/domain/user"
	"github.com/aburizalpurnama/travel/internal/app/middleware"
	"github.com/gofiber/fiber/v2"
)

type Option struct {
	Logger      *slog.Logger
	UserHandler *user.UserHandler
}

func SetupRoutesV1(app *fiber.App, opt *Option) {
	api := app.Group("/api/v1")
	api.Use(middleware.RequestLogger(opt.Logger))

	// Rute untuk User
	users := api.Group("/users")
	users.Post("/", opt.UserHandler.CreateUser)
	users.Get("/", opt.UserHandler.GetUsers)
	users.Get("/:id", opt.UserHandler.GetUser)
	// Tambahkan rute untuk Update dan Delete

}
