package httphelper

import (
	"net/http"

	"github.com/aburizalpurnama/travel/internal/pkg/apperror"
)

// MapErrorToHTTPStatus mapping apperror code to a standard HTTP status code.
func MapErrorToHTTPStatus(code apperror.Code) int {
	switch code {
	case
		apperror.UserNotFound,
		apperror.ProductNotFound,
		apperror.BookingNotFound:
		return http.StatusNotFound

	case
		apperror.EmailExists,
		apperror.DuplicateEntry,
		apperror.BookingAlreadyConfirmed:
		return http.StatusConflict

	case
		apperror.Validation,
		apperror.BadRequest:
		return http.StatusBadRequest

	case
		apperror.Unauthenticated:
		return http.StatusUnauthorized

	case
		apperror.Unauthorized:
		return http.StatusForbidden

	default:
		return http.StatusInternalServerError
	}
}
