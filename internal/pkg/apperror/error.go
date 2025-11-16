package apperror

import "fmt"

// AppError represents a standardized application error.
type AppError struct {
	Code    Code           `json:"code"`
	Message string         `json:"message"` // Internal message for logging
	Err     error          `json:"-"`       // Original error, not exposed in JSON
	Details map[string]any `json:"details,omitempty"`
}

// Error satisfies the standard error interface, providing a formatted error string.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap allows error unwrapping to retrieve the underlying original error.
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError instance with the specified code, message, original error, and additional details.
func New(code Code, message string, err error, details map[string]any) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
		Details: details,
	}
}
