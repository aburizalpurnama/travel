package user

import (
	"strconv"

	"github.com/aburizalpurnama/travel/internal/app/contract"
	"github.com/aburizalpurnama/travel/internal/app/model"
	"github.com/aburizalpurnama/travel/internal/app/payload"
	"github.com/aburizalpurnama/travel/internal/pkg/paginator"
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := h.service.CreateUser(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	option := paginator.OffsetBasedOption{}
	filter := model.UserFilter{}

	if page, err := strconv.Atoi(c.Query("page")); err == nil {
		option.Page = &page
	}

	if size, err := strconv.Atoi(c.Query("size")); err == nil {
		option.Size = &size
	}

	if id, err := strconv.Atoi(c.Query("id")); err == nil {
		id := uint(id)
		filter.ID = &id
	}

	if uid := c.Query("uid"); uid != "" {
		filter.UID = &uid
	}

	if email := c.Query("email"); email != "" {
		filter.Email = &email
	}

	if isActive, err := strconv.ParseBool(c.Query("is_active")); err == nil {
		filter.IsActive = &isActive
	}

	if isVerified, err := strconv.ParseBool(c.Query("is_verified")); err == nil {
		filter.IsVerified = &isVerified
	}

	if role := c.Query("role"); role != "" {
		filter.Role = &role
	}

	users, pagination, err := h.service.GetAllUsers(c.Context(), payload.UserGetAllRequest{Option: option, Filter: filter})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(user)
}

// Implementasi untuk UpdateUser dan DeleteUser mengikuti pola yang sama...
