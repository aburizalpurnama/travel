package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Product adalah model GORM untuk tabel "core.products"
type Product struct {
	ID          uint           `gorm:"primaryKey;autoIncrement"`
	UID         string         `gorm:"type:uuid;default:gen_random_uuid()"`
	CreatedOn   *time.Time     `gorm:"default:CURRENT_TIMESTAMP"`
	CreatedBy   datatypes.JSON `gorm:"type:jsonb;not null"`
	ModifiedOn  *time.Time
	ModifiedBy  datatypes.JSON  `gorm:"type:jsonb"`
	DeletedOn   gorm.DeletedAt  `gorm:"index"` // GORM soft delete
	Name        string          `gorm:"type:varchar(255);not null"`
	Description *string         `gorm:"type:text"`
	Price       decimal.Decimal `gorm:"type:decimal(18,2)"` // Gunakan float64 atau shopspring/decimal
	IsActive    *bool           `gorm:"default:true"`
}

// TableName menentukan nama tabel lengkap termasuk skema untuk GORM
func (Product) TableName() string {
	return "core.products"
}

type ProductFilter struct {
	IsActive *bool   `query:"is_active"`
	Search   *string `query:"search" search:"name,description"`
}
