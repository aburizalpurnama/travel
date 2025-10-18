package product

import (
	"strconv"

	"github.com/aburizalpurnama/travel/internal/app/contract"
	"github.com/aburizalpurnama/travel/internal/app/payload"
	"github.com/aburizalpurnama/travel/internal/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	service contract.ProductService
}

// NewProductHandler membuat instance baru dari ProductHandler
func NewProductHandler(service contract.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req payload.ProductCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ValidationError(err))
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ValidationError(err))
	}

	product, err := h.service.CreateProduct(c.Context(), req)
	if err != nil {
		c.Locals("error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error("internal", err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(product, nil))
}

func (h *ProductHandler) GetProducts(c *fiber.Ctx) error {
	req := payload.ProductGetAllRequest{}
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	req.SetDefault()

	products, pagination, err := h.service.GetAllProducts(c.Context(), req)
	if err != nil {
		c.Locals("error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error("internal", err.Error()))
	}

	return c.JSON(response.Success(products, pagination))
}

func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	product, err := h.service.GetProductByID(c.Context(), uint(id))
	if err != nil {
		c.Locals("error", err)
		return c.Status(fiber.StatusNotFound).JSON(response.Error("internal", err.Error()))
	}
	return c.JSON(response.Success(product, nil))
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var req payload.ProductUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ValidationError(err))
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ValidationError(err))
	}

	product, err := h.service.UpdateProduct(c.Context(), uint(id), req)
	if err != nil {
		c.Locals("error", err)
		return c.Status(fiber.StatusNotFound).JSON(response.Error("internal", err.Error()))
	}
	
	return c.JSON(response.Success(product, nil))
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	err = h.service.DeleteProduct(c.Context(), uint(id))
	if err != nil {
		c.Locals("error", err)
		return c.Status(fiber.StatusNotFound).JSON(response.Error("internal", err.Error()))
	}

	return c.JSON(response.Success("success delete data", nil))
}
