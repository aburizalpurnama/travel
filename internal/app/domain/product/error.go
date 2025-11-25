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
