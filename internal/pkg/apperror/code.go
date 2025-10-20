package apperror

type Code string

// GENERAL ERRORS
const (
	Unknown         Code = "ERR_UNKNOWN"
	Validation      Code = "ERR_VALIDATION"
	Internal        Code = "ERR_INTERNAL"
	NotFound        Code = "ERR_NOT_FOUND"
	Unauthenticated Code = "ERR_UNAUTHENTICATED"
	Unauthorized    Code = "ERR_UNAUTHORIZED"
	DuplicateEntry  Code = "ERR_DUPLICATE_ENTRY"
)

// DOMAIN ERRORS
const (
	// User
	UserNotFound    Code = "ERR_USER_NOT_FOUND"
	UserEmailExists Code = "ERR_USER_EMAIL_EXISTS"
	UserPhoneExists Code = "ERR_USER_PHONE_EXISTS"

	// Product
	ProductNotFound Code = "ERR_PRODUCT_NOT_FOUND"

	// TODO: define semua based on unique constraints
	ProductNameExists Code = "ERR_PRODUCT_NAME_EXISTS"

	// Finance
	InsufficientFunds Code = "ERR_INSUFFICIENT_FUNDS"
)

// VALIDATION DETAILS
const (
	IsRequired    Code = "IS_REQUIRED"
	IsInvalid     Code = "IS_INVALID"
	InvalidFormat Code = "INVALID_FORMAT"
	InvalidValue  Code = "INVALID_VALUE"
)
