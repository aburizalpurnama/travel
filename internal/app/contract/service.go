package contract

import (
	"context"

	"github.com/aburizalpurnama/travel/internal/app/payload"
	"github.com/aburizalpurnama/travel/internal/pkg/response"
)

// ProductService defines the business logic operations available for the Product model.
type ProductService interface {
	// CreateProduct handles the creation of a new product based on the provided request.
	CreateProduct(ctx context.Context, req payload.ProductCreateRequest) (*payload.ProductBaseResponse, error)

	// GetAllProducts retrieves a list of products matching the criteria in the request, including pagination.
	GetAllProducts(ctx context.Context, req payload.ProductGetAllRequest) ([]payload.ProductBaseResponse, *response.Pagination, error)

	// GetProductByID retrieves the details of a specific product identified by its ID.
	GetProductByID(ctx context.Context, id uint) (*payload.ProductBaseResponse, error)

	// UpdateProduct modifies an existing product identified by its ID with the provided update data.
	UpdateProduct(ctx context.Context, id uint, req payload.ProductUpdateRequest) (*payload.ProductBaseResponse, error)

	// DeleteProduct removes a product identified by its ID from the system.
	DeleteProduct(ctx context.Context, id uint) error
}
