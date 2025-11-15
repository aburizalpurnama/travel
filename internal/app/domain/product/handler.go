package product

import (
	"errors"
	"strconv"

	"github.com/aburizalpurnama/travel/internal/app/contract"
	"github.com/aburizalpurnama/travel/internal/app/payload"
	"github.com/aburizalpurnama/travel/internal/pkg/apperror"
	"github.com/aburizalpurnama/travel/internal/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var handlerTracer trace.Tracer = otel.Tracer("product.handler")

type Handler struct {
	service contract.ProductService
}

// NewHandler membuat instance baru dari ProductHandler
func NewHandler(service contract.ProductService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateProduct(c *fiber.Ctx) error {
	ctx, span := handlerTracer.Start(c.Context(), "CreateProduct")
	defer span.End()

	var req payload.ProductCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.JSONParserError(err))
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ValidationError(err))
	}

	product, err := h.service.CreateProduct(ctx, req)
	if err != nil {
		c.Locals("error", err)

		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			switch appErr.Code {
			case apperror.DuplicateEntry:
				return c.Status(fiber.StatusConflict).JSON(
					response.Error(appErr.Code, appErr.Message, appErr.Details),
				)

				// handle other error codes as needed

			default:
				return c.Status(fiber.StatusInternalServerError).JSON(
					response.Error(appErr.Code, appErr.Message, appErr.Details),
				)
			}
		}

		return c.Status(fiber.StatusInternalServerError).JSON(
			response.Error(apperror.Internal, apperror.ERR_INTERNAL_MSG, nil),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(product, nil))
}

func (h *Handler) GetProducts(c *fiber.Ctx) error {
	ctx, span := handlerTracer.Start(c.Context(), "GetProducts")
	defer span.End()

	req := payload.ProductGetAllRequest{}
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	req.SetDefault()

	products, pagination, err := h.service.GetAllProducts(ctx, req)
	if err != nil {
		c.Locals("error", err)

		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			switch appErr.Code {
			case apperror.NotFound:
				return c.Status(fiber.StatusConflict).JSON(
					response.Error(appErr.Code, appErr.Message, appErr.Details),
				)

				// handle other error codes as needed

			default:
				return c.Status(fiber.StatusInternalServerError).JSON(
					response.Error(appErr.Code, appErr.Message, appErr.Details),
				)
			}
		}

		return c.Status(fiber.StatusInternalServerError).JSON(
			response.Error(apperror.Internal, apperror.ERR_INTERNAL_MSG, nil),
		)
	}

	return c.JSON(response.Success(products, pagination))
}

func (h *Handler) GetProduct(c *fiber.Ctx) error {
	ctx, span := handlerTracer.Start(c.Context(), "GetProduct")
	defer span.End()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.Error(apperror.Validation, "invalid id", nil),
		)
	}

	product, err := h.service.GetProductByID(ctx, uint(id))
	if err != nil {
		c.Locals("error", err)

		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			switch appErr.Code {
			case apperror.ProductNotFound:
				return c.Status(fiber.StatusNotFound).JSON(
					response.Error(appErr.Code, appErr.Message, appErr.Details),
				)

				// handle other error codes as needed
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(
					response.Error(apperror.Internal, apperror.ERR_INTERNAL_MSG, nil),
				)
			}
		}

		return c.Status(fiber.StatusInternalServerError).JSON(
			response.Error(apperror.Internal, apperror.ERR_INTERNAL_MSG, nil),
		)
	}

	return c.JSON(response.Success(product, nil))
}

func (h *Handler) UpdateProduct(c *fiber.Ctx) error {
	ctx, span := handlerTracer.Start(c.Context(), "UpdateProduct")
	defer span.End()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var req payload.ProductUpdateRequest
	err = c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.JSONParserError(err))
	}

	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ValidationError(err))
	}

	product, err := h.service.UpdateProduct(ctx, uint(id), req)
	if err != nil {
		c.Locals("error", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			response.Error(apperror.Internal, apperror.ERR_INTERNAL_MSG, nil),
		)
	}

	return c.JSON(response.Success(product, nil))
}

func (h *Handler) DeleteProduct(c *fiber.Ctx) error {
	ctx, span := handlerTracer.Start(c.Context(), "DeleteProduct")
	defer span.End()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	err = h.service.DeleteProduct(ctx, uint(id))
	if err != nil {
		c.Locals("error", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			response.Error(apperror.Internal, apperror.ERR_INTERNAL_MSG, nil),
		)
	}

	return c.JSON(response.Success("success delete data", nil))
}
