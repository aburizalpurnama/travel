package response

import "math"

type APIResponse struct {
	Status     string      `json:"status"`
	Data       any         `json:"data,omitempty"`
	Error      *APIError   `json:"error,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
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

func Error(code, message string) APIResponse {
	return APIResponse{
		Status: "error",
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	}
}

func ValidationError(details any) APIResponse {
	return APIResponse{
		Status: "error",
		Error: &APIError{
			Code:    "VALIDATION_ERROR",
			Message: "Your request is invalid. Please check the details.",
			Details: details,
		},
	}
}
