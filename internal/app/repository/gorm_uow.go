package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/aburizalpurnama/travel/internal/app/contract"
	"github.com/aburizalpurnama/travel/internal/app/domain/product"
)

// gormUnitOfWork implements the contract.UnitOfWork interface.
// It acts as the main provider for repositories and transaction management.
type gormUnitOfWork struct {
	db *gorm.DB // Main database connection pool

	// Caches for lazy-loaded repositories
	productRepo contract.ProductRepository
}

// NewGormUnitOfWork creates a new UnitOfWork provider.
func NewGormUnitOfWork(db *gorm.DB) contract.UnitOfWork {
	return &gormUnitOfWork{db: db}
}

// Product provides a lazy-loaded transactional ProductRepository.
func (u *gormUnitOfWork) Product() contract.ProductRepository {
	if u.productRepo == nil {
		u.productRepo = product.NewRepository(u.db)
	}
	return u.productRepo
}

// Execute runs the given function 'fn' within a single GORM transaction.
// If 'fn' returns an error, GORM automatically performs a rollback.
// If 'fn' succeeds, GORM automatically performs a commit.
func (u *gormUnitOfWork) Execute(ctx context.Context, fn func(contract.UnitOfWork) error) error {
	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create a new transactional UoW instance bound to this specific tx
		txUoW := &gormUnitOfWork{db: tx}

		// Pass the transactional UoW to the business logic function
		return fn(txUoW)
	})
}
