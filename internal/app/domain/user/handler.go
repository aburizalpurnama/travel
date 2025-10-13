package user

import (
	"strconv"

	"github.com/aburizalpurnama/travel/internal/app/contract"
	"github.com/aburizalpurnama/travel/internal/app/payload"
	"github.com/aburizalpurnama/travel/internal/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service contract.UserService
}

// NewUserHandler membuat instance baru dari UserHandler
func NewUserHandler(service contract.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req payload.UserCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ValidationError(err))
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ValidationError(err))
	}

	user, err := h.service.CreateUser(c.Context(), req)
	if err != nil {
		c.Locals("error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error("internal", err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(user, nil))
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	req := payload.UserGetAllRequest{}
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	req.SetDefault()

	users, pagination, err := h.service.GetAllUsers(c.Context(), req)
	if err != nil {
		c.Locals("error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error("internal", err.Error()))
	}

	return c.JSON(response.Success(users, pagination))
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	user, err := h.service.GetUserByID(c.Context(), uint(id))
	if err != nil {
		c.Locals("error", err)
		return c.Status(fiber.StatusNotFound).JSON(response.Error("internal", err.Error()))
	}
	return c.JSON(user)
}

// Implementasi untuk UpdateUser dan DeleteUser mengikuti pola yang sama...
