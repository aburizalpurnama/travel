package model

import (
	"time"

	"gorm.io/datatypes"
)

// User merepresentasikan model data di database
type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	UID          string    `gorm:"type:uuid;default:gen_random_uuid()"`
	CreatedOn    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	CreatedBy    datatypes.JSON
	ModifiedOn   *time.Time
	ModifiedBy   datatypes.JSON
	DeletedOn    *time.Time
	FirstName    string  `gorm:"type:varchar(100);not null"`
	MiddleName   *string `gorm:"type:varchar(100)"`
	LastName     *string `gorm:"type:varchar(100)"`
	FullName     string  `gorm:"type:varchar(255);not null"`
	Gender       string  `gorm:"type:user.customer_gender_enum;not null"`
	Email        *string `gorm:"type:varchar(320);uniqueIndex:ux_users_email_active,where:deleted_on IS NULL"`
	Phone        string  `gorm:"type:varchar(50);not null;uniqueIndex:ux_users_phone_active,where:deleted_on IS NULL"`
	PasswordHash *string `gorm:"type:varchar(255)"`
	IsActive     bool    `gorm:"default:true"`
	IsVerified   bool    `gorm:"default:false"`
	Role         string  `gorm:"type:user.users_role_enum;not null"`
}

func (User) TableName() string {
	return "user.users"
}

type UserFilter struct {
	UID        *string `query:"uid"`
	IsActive   *bool   `query:"is_active"`
	IsVerified *bool   `query:"is_verified"`
	Role       *string `query:"role"`
	Search     *string `query:"search" search:"full_name,email,phone"`
}
