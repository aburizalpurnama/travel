package product

import (
	"github.com/aburizalpurnama/travel/internal/app/contract"
	"github.com/aburizalpurnama/travel/internal/app/model"
	"github.com/aburizalpurnama/travel/internal/pkg/repository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// repositoryTracer
var _ trace.Tracer = otel.Tracer("product.repository")

// Repository implements the contract.ProductRepository interface.
// It embeds a generic GORM repository to handle basic CRUD operations.
type Repository struct {
	*repository.GORM[model.Product, model.ProductFilter]
	db *gorm.DB
}

// NewRepository creates a new product repository instance.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		GORM: repository.NewGORM[model.Product, model.ProductFilter](db),
		db:   db,
	}
}

// Ensures implementaton satisfies the contract at compile-time.
var _ contract.ProductRepository = (*Repository)(nil)

// Add your custom repository methods below
