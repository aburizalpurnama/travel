package mapper

import (
	"github.com/aburizalpurnama/travel/internal/app/contract"
	"github.com/jinzhu/copier"
)

// copierMapper implements the contract.Mapper interface using the 'copier' library.
type copierMapper struct {
	// options can be added here
}

// NewCopierMapper creates a new mapper instance.
func NewCopierMapper() contract.Mapper {
	return &copierMapper{}
}

// Ensures implementation satisfies the contract at compile-time.
var _ contract.Mapper = (*copierMapper)(nil)

// ToModel maps a Request DTO to a Domain Model, ignoring nil fields from the DTO.
// This is ideal for 'Update' operations.
func (m *copierMapper) ToModel(dto any, model any) error {
	return copier.CopyWithOption(model, dto, copier.Option{
		IgnoreEmpty: true, // only copy not nil/zero-value
	})
}

// ToResponse maps a Domain Model to a Response DTO.
func (m *copierMapper) ToResponse(model any, response any) error {
	return copier.Copy(response, model)
}
