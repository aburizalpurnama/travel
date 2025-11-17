package product

import "github.com/aburizalpurnama/travel/internal/pkg/apperror"

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
