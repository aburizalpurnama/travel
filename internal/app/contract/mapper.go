package contract

// Mapper defines the contract for model mapping.
type Mapper interface {
	// ToModel maps a Request DTO (src) to a Domain Model (dst).
	// 'dst' must be a pointer.
	ToModel(src any, dst any) error

	// ToResponse maps a Domain Model (src) to a Response DTO (dst).
	// 'dst' must be a pointer.
	ToResponse(src any, dst any) error
}
