package payload

import (
	"time"

	"github.com/aburizalpurnama/travel/internal/app/model"
)

// ==========================================================
// DTO untuk REQUEST (Data Masuk)
// ==========================================================

type UserGetAllRequest struct {
	*CommonGetAllRequest
	*model.UserFilter
}

type UserCreateRequest struct {
	FirstName            string `json:"first_name" validate:"required,max=100"`
	MiddleName           string `json:"middle_name,omitempty" validate:"omitempty,max=100"`
	LastName             string `json:"last_name,omitempty" validate:"omitempty,max=100"`
	FullName             string `json:"full_name" validate:"required,max=255"`
	Gender               string `json:"gender" validate:"required,oneof=male female"`
	Email                string `json:"email" validate:"required,email,max=320"`
	Phone                string `json:"phone" validate:"required,max=50"`
	Password             string `json:"password" validate:"required,min=8"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
	Role                 string `json:"role" validate:"required,oneof=customer muthawif"`
}

type UserUpdateRequest struct {
	FirstName  *string `json:"first_name,omitempty" validate:"omitempty,max=100"`
	MiddleName *string `json:"middle_name,omitempty" validate:"omitempty,max=100"`
	LastName   *string `json:"last_name,omitempty" validate:"omitempty,max=100"`
	FullName   *string `json:"full_name,omitempty" validate:"omitempty,max=255"`
	Gender     *string `json:"gender,omitempty" validate:"omitempty,oneof=male female"`
	Phone      *string `json:"phone,omitempty" validate:"omitempty,max=50"`
}

// ==========================================================
// DTO untuk RESPONSE (Data Keluar)
// ==========================================================

type UserResponse struct {
	UID        string    `json:"uid"`
	FirstName  string    `json:"first_name"`
	MiddleName *string   `json:"middle_name,omitempty"`
	LastName   *string   `json:"last_name,omitempty"`
	FullName   string    `json:"full_name"`
	Gender     string    `json:"gender"`
	Email      *string   `json:"email,omitempty"`
	Phone      string    `json:"phone"`
	IsActive   bool      `json:"is_active"`
	IsVerified bool      `json:"is_verified"`
	Role       string    `json:"role"`
	CreatedOn  time.Time `json:"created_on"`
}
