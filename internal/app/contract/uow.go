package contract

import (
	"context"
)

// UnitOfWork defines the interface for managing atomic database operations (transactions) and provides access to repositories.
type UnitOfWork interface {
	ProductRepository() ProductRepository

	// RunInTransaction runs the given function 'fn' within a single atomic transaction.
	// If 'fn' returns an error, the transaction is rolled back.
	// If 'fn' succeeds, the transaction is committed.
	RunInTransaction(ctx context.Context, fn func(context.Context, UnitOfWork) error) error
}
