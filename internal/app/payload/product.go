package payload

import (
	"time"

	"github.com/aburizalpurnama/travel/internal/app/model"
)

// ==========================================================
// Request DTOs
// ==========================================================

// ProductGetAllRequest defines the query parameters for retrieving a list of products.
// It combines common pagination/sorting parameters with specific product filters.
type ProductGetAllRequest struct {
	*CommonGetAllRequest
	*model.ProductFilter
}

// ProductCreateRequest defines the payload required to create a new product.
type ProductCreateRequest struct {
	Name        string  `json:"name" validate:"required,max=255"`
	Description *string `json:"description,omitempty"`
	Price       string  `json:"price,omitempty" validate:"omitempty,gt=0"`
	IsActive    *bool   `json:"is_active,omitempty" validate:"omitempty"`
}

// ProductUpdateRequest defines the payload for updating an existing product.
// All fields are optional to allow partial updates.
type ProductUpdateRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,max=255"`
	Description *string `json:"description,omitempty"`
	Price       string  `json:"price,omitempty" validate:"omitempty,gt=0"`
	IsActive    *bool   `json:"is_active,omitempty" validate:"omitempty"`
}

// ==========================================================
// Response DTOs
// ==========================================================

// ProductBaseResponse defines the standard response structure for product data.
type ProductBaseResponse struct {
	ID          uint      `json:"id"`
	UID         string    `json:"uid"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Price       string    `json:"price,omitempty"`
	IsActive    *bool     `json:"is_active"`
	CreatedOn   time.Time `json:"created_on"`
}
