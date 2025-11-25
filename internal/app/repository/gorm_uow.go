package repository

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"

	"github.com/aburizalpurnama/travel/internal/app/contract"
	"github.com/aburizalpurnama/travel/internal/app/domain/product"
)

var tracer trace.Tracer = otel.Tracer("repository.uow")

// gormUnitOfWork implements the contract.UnitOfWork interface.
// It acts as the main provider for repositories and transaction management.
type gormUnitOfWork struct {
	db *gorm.DB // Main database connection pool

	// Caches for lazy-loaded repositories
	productRepo contract.ProductRepository
}

// NewGORMUnitOfWork creates a new UnitOfWork provider with GORM DB.
func NewGORMUnitOfWork(db *gorm.DB) contract.UnitOfWork {
	return &gormUnitOfWork{db: db}
}

// Ensures implementation satisfies the contract at compile-time.
var _ contract.UnitOfWork = (*gormUnitOfWork)(nil)

// ProductRepository provides a lazy-loaded transactional ProductRepository.
func (u *gormUnitOfWork) ProductRepository() contract.ProductRepository {
	if u.productRepo == nil {
		u.productRepo = product.NewRepository(u.db)
	}
	return u.productRepo
}

// RunInTransaction runs the given function 'fn' within a single GORM transaction.
// If 'fn' returns an error, GORM automatically performs a rollback.
// If 'fn' succeeds, GORM automatically performs a commit.
func (u *gormUnitOfWork) RunInTransaction(ctx context.Context, fn func(context.Context, contract.UnitOfWork) error) error {
	ctx, span := tracer.Start(ctx, "Execute") // Span untuk seluruh transaksi
	defer span.End()

	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create a new transactional UoW instance bound to this specific tx
		txUoW := &gormUnitOfWork{db: tx}

		// Pass the transactional UoW to the business logic function
		return fn(ctx, txUoW)
	})
}
