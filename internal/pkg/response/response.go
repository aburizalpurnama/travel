package response

type APIResponse[T any] struct {
	Status     string      `json:"status"`
	Data       T           `json:"data,omitempty"`
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

func Success[T any](data T, pagination *Pagination) APIResponse[T] {
	return APIResponse[T]{
		Status:     "success",
		Data:       data,
		Pagination: pagination,
	}
}

func Error[T any](code, message string) APIResponse[T] {
	return APIResponse[T]{
		Status: "error",
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	}
}

func ValidationError[T any](details any) APIResponse[T] {
	return APIResponse[T]{
		Status: "error",
		Error: &APIError{
			Code:    "VALIDATION_ERROR",
			Message: "Your request is invalid. Please check the details.",
			Details: details,
		},
	}
}
