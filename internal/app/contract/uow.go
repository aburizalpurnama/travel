package contract

import (
	"context"
)

// UnitOfWork defines the interface for managing atomic database operations (transactions).
type UnitOfWork interface {
	// Product returns a ProductRepository bound to this unit of work.
	Product() ProductRepository

	// Execute runs the given function 'fn' within a single atomic transaction.
	// If 'fn' returns an error, the transaction is rolled back.
	// If 'fn' succeeds, the transaction is committed.
	Execute(ctx context.Context, fn func(context.Context, UnitOfWork) error) error
}
