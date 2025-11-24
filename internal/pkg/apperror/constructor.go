package apperror

// ==========================================================
// General Error Constructors
// ==========================================================

// ErrMapping creates a new error for data mapping failures.
func ErrMapping(err error) *AppError {
	return New(
		Internal,
		"failed to map data",
		err,
		nil,
	)
}
