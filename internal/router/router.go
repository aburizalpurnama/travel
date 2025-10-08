package router

import (
	"github.com/aburizalpurnama/travel/internal/user" // Pastikan ini benar
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, userHandler *user.UserHandler) {
	api := app.Group("/api/v1")

	// Rute untuk User
	users := api.Group("/users")
	users.Post("/", userHandler.CreateUser)
	users.Get("/", userHandler.GetUsers)
	users.Get("/:id", userHandler.GetUser)
	// Tambahkan rute untuk Update dan Delete
}
