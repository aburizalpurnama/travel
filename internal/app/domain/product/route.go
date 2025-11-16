package product

import "github.com/gofiber/fiber/v2"

// NewRoute registers product-related routes to the provided router group.
func NewRoute(router fiber.Router, handler *Handler) {
	products := router.Group("/products")

	products.Post("/", handler.CreateProduct)
	products.Get("/", handler.GetProducts)
	products.Get("/:id", handler.GetProduct)
	products.Patch("/:id", handler.UpdateProduct)
	products.Delete("/:id", handler.DeleteProduct)
}
