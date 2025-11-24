package product

import "github.com/aburizalpurnama/travel/internal/pkg/apperror"

// ==========================================================
// Product Error Constructors
// ==========================================================

// ErrProductNotFound creates a new error for missing product records.
func ErrProductNotFound(err error) *apperror.AppError {
	return apperror.New(
		apperror.ProductNotFound,
		"product not found",
		err,
		nil,
	)
}

// ErrFailedCreateProduct creates a new error for failed create product.
func ErrFailedCreateProduct(err error) *apperror.AppError {
	return apperror.New(
		apperror.Internal,
		"failed to create product",
		err,
		nil,
	)
}

// ErrFailedUpdateProduct creates a new error for failures during product updates.
func ErrFailedUpdateProduct(err error) *apperror.AppError {
	return apperror.New(
		apperror.Internal,
		"failed to update product",
		err,
		nil,
	)
}

// ErrFailedDeleteProduct creates a new error for failures during product deletion.
func ErrFailedDeleteProduct(err error) *apperror.AppError {
	return apperror.New(
		apperror.Internal,
		"failed to delete product",
		err,
		nil,
	)
}
