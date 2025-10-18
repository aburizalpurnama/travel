package payload

import (
	"time"

	"github.com/aburizalpurnama/travel/internal/app/model"
)

// ==========================================================
// REQUEST DTOs (Data Masuk dari Klien)
// ==========================================================

type ProductGetAllRequest struct {
	*CommonGetAllRequest
	*model.ProductFilter
}

type ProductCreateRequest struct {
	Name        string  `json:"name" validate:"required,max=255"`
	Description *string `json:"description,omitempty"`
	Price       string  `json:"price,omitempty" validate:"omitempty,gt=0"`
	IsActive    *bool   `json:"is_active,omitempty" validate:"omitempty"`
}

type ProductUpdateRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,max=255"`
	Description *string `json:"description,omitempty"`
	Price       string  `json:"price,omitempty" validate:"omitempty,gt=0"`
	IsActive    *bool   `json:"is_active,omitempty" validate:"omitempty"`
}

// ==========================================================
// RESPONSE DTO (Data Keluar ke Klien)
// ==========================================================

type ProductBaseResponse struct {
	ID          uint      `json:"id"`
	UID         string    `json:"uid"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Price       string    `json:"price,omitempty"`
	IsActive    *bool     `json:"is_active"`
	CreatedOn   time.Time `json:"created_on"`
}
