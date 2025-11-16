package contract

import (
	"context"

	"github.com/aburizalpurnama/travel/internal/app/model"
)

// ProductRepository defines the standard database operations for the Product model.
type ProductRepository interface {
	// FindAll retrieves a list of products based on pagination parameters and filter criteria.
	FindAll(ctx context.Context, page *int, size *int, filter *model.ProductFilter) ([]model.Product, error)

	// Count returns the total number of products that match the given filter.
	Count(ctx context.Context, filter *model.ProductFilter) (int64, error)

	// FindByID retrieves a single product by its unique identifier.
	FindByID(ctx context.Context, id uint) (*model.Product, error)

	// Save persists a new product record to the database.
	Save(ctx context.Context, product *model.Product) (*model.Product, error)

	// Update modifies an existing product record in the database.
	Update(ctx context.Context, product *model.Product) (*model.Product, error)

	// Delete removes a product record from the database by its ID.
	Delete(ctx context.Context, id uint) error
}
