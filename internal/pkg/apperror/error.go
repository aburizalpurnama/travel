package apperror

import "fmt"

type AppError struct {
	Code    Code           `json:"code"`
	Message string         `json:"message"` // Pesan internal untuk logging
	Err     error          `json:"-"`       // Error asli, jangan diekspos ke JSON
	Details map[string]any `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code Code, message string, err error, details map[string]any) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
		Details: details,
	}
}
