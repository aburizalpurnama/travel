package response

import (
	"encoding/json"
	"errors"
	"math"

	"github.com/aburizalpurnama/travel/internal/pkg/apperror"
	appstrings "github.com/aburizalpurnama/travel/internal/pkg/strings"
	"github.com/go-playground/validator/v10"
)

// APIResponse represents the standard JSON response structure for the API.
type APIResponse struct {
	Status     string      `json:"status"`
	Data       any         `json:"data,omitempty"`
	Error      *APIError   `json:"error,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// APIError represents the structure of error details in the response.
type APIError struct {
	Code    apperror.Code `json:"code"`
	Message string        `json:"message"`
	Details any           `json:"details,omitempty"`
}

// Pagination holds pagination metadata.
type Pagination struct {
	TotalItems  *int64 `json:"total_items,omitempty"`
	TotalPages  *int   `json:"total_pages,omitempty"`
	CurrentPage *int   `json:"current_page,omitempty"`
	PageSize    *int   `json:"page_size,omitempty"`
}

// NewPagination creates a new Pagination instance based on page size and total count.
func NewPagination(page, size *int, count *int64) *Pagination {
	if page == nil || size == nil || count == nil {
		return nil
	}

	if *size <= 0 || *count <= 0 {
		return nil
	}

	totalPages := int(math.Ceil(float64(*count) / float64(*size)))
	return &Pagination{
		CurrentPage: page,
		PageSize:    size,
		TotalItems:  count,
		TotalPages:  &totalPages,
	}
}

// Success creates a success APIResponse with data and optional pagination.
func Success(data any, pagination *Pagination) APIResponse {
	return APIResponse{
		Status:     "success",
		Data:       data,
		Pagination: pagination,
	}
}

// Error creates an error APIResponse with a code, message, and optional details.
func Error(code apperror.Code, message string, details any) APIResponse {
	return APIResponse{
		Status: "error",
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

// JSONParserError converts JSON unmarshalling errors into a standardized APIResponse.
func JSONParserError(err error) APIResponse {
	var unmarshalErr *json.UnmarshalTypeError
	var syntaxErr *json.SyntaxError

	if errors.As(err, &unmarshalErr) {
		return APIResponse{
			Status: "error",
			Error: &APIError{
				Code:    apperror.Validation,
				Message: "Your request is invalid. Please check the details.",
				Details: map[string]any{unmarshalErr.Field: apperror.InvalidFormat},
			},
		}
	}

	if errors.As(err, &syntaxErr) {
		return APIResponse{
			Status: "error",
			Error: &APIError{
				Code:    apperror.Validation,
				Message: "Your request body is malformed JSON.",
			},
		}
	}

	return APIResponse{
		Status: "error",
		Error: &APIError{
			Code:    apperror.Validation,
			Message: "Your request body is invalid.",
		},
	}
}

// QueryParserError converts query parsing errors into a standardized APIResponse.
func QueryParserError(err error) APIResponse {
  return APIResponse{
    Status: "error",
    Error: &APIError{
      Code:    apperror.Validation,
      Message: "Your query parameters are invalid.",
      Details: map[string]any{
        "syntax": err.Error(),
      },
    },
  }
}

// ValidationError converts go-playground/validator errors into a standardized APIResponse.
// It maps validator tags to custom apperror codes.
func ValidationError(err error) APIResponse {
	formattedErrors := make(map[string]apperror.Code)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range validationErrs {
			field := appstrings.ToSnakeCase(fe.Field())

			// Map validator tags to custom apperror codes
			switch fe.Tag() {
			case "required":
				formattedErrors[field] = apperror.IsRequired
			case "email", "url", "uuid":
				formattedErrors[field] = apperror.InvalidFormat
			case "gt", "gte":
				formattedErrors[field] = apperror.ValueTooLow
			case "lt", "lte":
				formattedErrors[field] = apperror.ValueTooHigh
			case "min":
				formattedErrors[field] = apperror.LengthTooShort
			case "max":
				formattedErrors[field] = apperror.LengthTooLong
			case "length":
				formattedErrors[field] = apperror.InvalidLength
			case "oneof":
				formattedErrors[field] = apperror.InvalidChoice

			default:
				formattedErrors[field] = apperror.InvalidValue
			}
		}

		return APIResponse{
			Status: "error",
			Error: &APIError{
				Code:    apperror.Validation,
				Message: "Your request is invalid. Please check the details.",
				Details: formattedErrors,
			},
		}
	}

	return APIResponse{
		Status: "error",
		Error: &APIError{
			Code:    apperror.Validation,
			Message: "Your request is invalid",
		},
	}
}
