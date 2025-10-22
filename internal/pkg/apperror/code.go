package apperror

// Code is a custom type for stable, machine-readable error codes.
// These codes are intended to be consumed by clients (frontend/mobile)
// to perform translation or take specific actions.
type Code string

// ==========================================================
// GENERAL ERROR CODES (Typically mapped from HTTP Status Codes)
// ==========================================================
const (
	// ERR_UNKNOWN (500)
	// Fallback for completely unexpected errors.
	Unknown Code = "ERR_UNKNOWN"

	// ERR_INTERNAL (500)
	// Known server-side errors (e.g., failed to connect to DB, recovered panic).
	Internal Code = "ERR_INTERNAL"

	// ERR_BAD_REQUEST (400)
	// General error for invalid client requests (e.g., malformed JSON).
	BadRequest Code = "ERR_BAD_REQUEST"

	// ERR_VALIDATION (400)
	// Specific error for when DTO data validation fails.
	// This is usually accompanied by a 'details' field.
	Validation Code = "ERR_VALIDATION"

	// ERR_UNAUTHENTICATED (401)
	// The client tried to access a protected resource without a valid token.
	Unauthenticated Code = "ERR_UNAUTHENTICATED"

	// ERR_TOKEN_EXPIRED (401)
	// The provided authentication token has expired.
	TokenExpired Code = "ERR_TOKEN_EXPIRED"

	// ERR_UNAUTHORIZED (403)
	// The client is authenticated (logged in) but lacks the necessary permissions (role)
	// to access this specific resource.
	Unauthorized Code = "ERR_UNAUTHORIZED"

	// ERR_NOT_FOUND (404)
	// A fallback code if a more domain-specific 'Not Found' error isn't used.
	NotFound Code = "ERR_NOT_FOUND"

	// ERR_DUPLICATE_ENTRY (409 Conflict)
	// General code for a unique constraint violation in the database.
	DuplicateEntry Code = "ERR_DUPLICATE_ENTRY"

	// ERR_STATE_CONFLICT (409 Conflict)
	// The action cannot be performed due to the current state of the resource.
	// (e.g., trying to edit a 'canceled' booking).
	StateConflict Code = "ERR_STATE_CONFLICT"

	// ERR_RATE_LIMIT_EXCEEDED (429)
	// The client has sent too many requests in a given amount of time.
	RateLimitExceeded Code = "ERR_RATE_LIMIT_EXCEEDED"

	// ERR_SERVICE_UNAVAILABLE (503)
	// An error caused by a failing external (third-party) service.
	// (e.g., the payment gateway API is down).
	ServiceUnavailable Code = "ERR_SERVICE_UNAVAILABLE"
)

// ==========================================================
// DOMAIN-SPECIFIC ERROR CODES (Business Logic)
// ==========================================================
const (
	// User (ERR_USER_...)
	UserNotFound   Code = "ERR_USER_NOT_FOUND"
	EmailExists    Code = "ERR_EMAIL_EXISTS"
	PhoneExists    Code = "ERR_PHONE_EXISTS"
	ReferralExists Code = "ERR_REFERRAL_CODE_EXISTS"

	// Product (ERR_PRODUCT_...)
	ProductNotFound   Code = "ERR_PRODUCT_NOT_FOUND"
	ProductNameExists Code = "ERR_PRODUCT_NAME_EXISTS"
	SKUExists         Code = "ERR_SKU_EXISTS"

	// Booking (ERR_BOOKING_...)
	BookingNotFound         Code = "ERR_BOOKING_NOT_FOUND"
	BookingAlreadyConfirmed Code = "ERR_BOOKING_ALREADY_CONFIRMED"
	BookingNotCancellable   Code = "ERR_BOOKING_NOT_CANCELLABLE"
	BookingBatchSoldOut     Code = "ERR_BOOKING_BATCH_SOLD_OUT"

	// Finance (ERR_FIN_...)
	InsufficientFunds Code = "ERR_INSUFFICIENT_FUNDS"
	VoucherExpired    Code = "ERR_VOUCHER_EXPIRED"
)

// ==========================================================
// VALIDATION DETAIL CODES (For the 'details' field in ERR_VALIDATION)
// ==========================================================
const (
	// INVALID_VALUE
	// The fallback error code (default if not match code found).
	InvalidValue Code = "INVALID_VALUE"

	// IS_REQUIRED
	// The field is missing, nil, or empty.
	IsRequired Code = "IS_REQUIRED"

	// INVALID_TYPE
	// The field's type is incorrect
	InvalidType Code = "INVALID_TYPE"

	// INVALID_FORMAT
	// The field's format is incorrect (e.g., email, url, uuid, date).
	InvalidFormat Code = "INVALID_FORMAT"

	// INVALID_LENGTH
	// The string or array length is invalid.
	InvalidLength Code = "INVALID_LENGTH"

	// LENGTH_TOO_SHORT
	// The string or array is shorter than the minimum allowed (min).
	LengthTooShort Code = "LENGTH_TOO_SHORT"

	// LENGTH_TOO_LONG
	// The string or array is longer than the maximum allowed (max).
	LengthTooLong Code = "LENGTH_TOO_LONG"

	// VALUE_TOO_LOW
	// The number is smaller than the minimum allowed (gt, gte).
	ValueTooLow Code = "VALUE_TOO_LOW"

	// VALUE_TOO_HIGH
	// The number is larger than the maximum allowed (lt, lte).
	ValueTooHigh Code = "VALUE_TOO_HIGH"

	// INVALID_CHOICE
	// The value is not one of the allowed enum options.
	InvalidChoice Code = "INVALID_CHOICE"
)
