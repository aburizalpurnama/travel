package response

import (
	"encoding/json"
	"errors"
	"math"

	"github.com/aburizalpurnama/travel/internal/pkg/apperror"
	"github.com/aburizalpurnama/travel/internal/pkg/dberror"
	appstrings "github.com/aburizalpurnama/travel/internal/pkg/strings"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
)

type APIResponse struct {
	Status     string      `json:"status"`
	Data       any         `json:"data,omitempty"`
	Error      *APIError   `json:"error,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type APIError struct {
	Code    apperror.Code `json:"code"`
	Message string        `json:"message"`
	Details any           `json:"details,omitempty"`
}

type Pagination struct {
	TotalItems  *int64 `json:"total_items,omitempty"`
	TotalPages  *int   `json:"total_pages,omitempty"`
	CurrentPage *int   `json:"current_page,omitempty"`
	PageSize    *int   `json:"page_size,omitempty"`
}

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

func Success(data any, pagination *Pagination) APIResponse {
	return APIResponse{
		Status:     "success",
		Data:       data,
		Pagination: pagination,
	}
}

func Error(code apperror.Code, message string) APIResponse {
	return APIResponse{
		Status: "error",
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	}
}

// func ValidationError(details any) APIResponse {
// 	return APIResponse{
// 		Status: "error",
// 		Error: &APIError{
// 			Code:    apperror.Validation,
// 			Message: "Your request is invalid. Please check the details.",
// 			Details: details,
// 		},
// 	}
// }

func JSONParserError(err error) APIResponse {
	var unmarshalErr *json.UnmarshalTypeError
	var syntaxErr *json.SyntaxError

	if errors.As(err, &unmarshalErr) {
		// expectedType := unmarshalErr.Type.String()

		// switch true {
		// case strings.Contains(expectedType, "float"), strings.Contains(expectedType, "int"):
		// 	expectedType = "number"
		// case strings.Contains(expectedType, "bool"):
		// 	expectedType = "boolean"
		// case strings.Contains(expectedType, "string"):
		// 	expectedType = "string"
		// default:
		// 	expectedType = "a valid value"
		// }

		// unmarchalErr.Value

		return APIResponse{
			Status: "error",
			Error: &APIError{
				Code:    apperror.Validation,
				Message: "Your request is invalid. Please check the details.",
				Details: map[string]any{unmarshalErr.Field: apperror.IsInvalid},
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

func ValidationError(err error) APIResponse {
	formattedErrors := make(map[string]apperror.Code)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range validationErrs {
			field := appstrings.ToSnakeCase(fe.Field())

			// Petakan 'tag' validasi ke 'apperror.Code' kustom Anda
			switch fe.Tag() {
			case "required":
				formattedErrors[field] = apperror.IsRequired
			case "email", "url", "uuid":
				formattedErrors[field] = apperror.InvalidFormat
			case "gt", "min", "max", "length":
				formattedErrors[field] = apperror.InvalidValue
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

func PostgreError(pgErr *pgconn.PgError) APIResponse {
	// if pgErr.Code == "23505" {

	// 			// TODO: define semua based on unique constraints
	// 			// (Opsional tapi sangat baik) Cek nama constraint-nya
	// 			switch pgErr.ConstraintName {
	// 			case "products_name_unique":
	// 				return nil, apperror.New(apperror.ProductNameExists, "product name already exists", err)
	// 			case "ux_products_uid_active":
	// 				return nil, apperror.New(apperror.DuplicateEntry, "product UID already exists", err)
	// 			default:
	// 				return nil, apperror.New(apperror.DuplicateEntry, "unique constraint violated", err)
	// 			}
	// 		}

	// GetSQLState

	if dberror.SQLState(pgErr.Code) == dberror.UniqueViolation {
		return APIResponse{
			Status: "error",
			Error: &APIError{
				Code: apperror.DuplicateEntry,
				/ dapetin error message dari pgErr.Detail
				Message: "Duplicate entry violates unique constraint",
			},
		}
	}

	return APIResponse{
		Status: "error",
		Error: &APIError{
			Code:    apperror.Internal,
			Message: apperror.ERR_INTERNAL_MSG,
		},
	}
}
