package domain

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

// UserCreateRequest adalah DTO untuk membuat user baru
type UserCreateRequest struct {
	FirstName string `json:"first_name" validate:"required,max=100"`
	FullName  string `json:"full_name" validate:"required,max=255"`
	Gender    string `json:"gender" validate:"required,oneof=male female"`
	Email     string `json:"email" validate:"required,email,max=320"`
	Phone     string `json:"phone" validate:"required,max=50"`
	Password  string `json:"password" validate:"required,min=8"`
	Role      string `json:"role" validate:"required,oneof=customer muthawif"`
}

// UserUpdateRequest adalah DTO untuk memperbarui user
type UserUpdateRequest struct {
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,max=100"`
	FullName  *string `json:"full_name,omitempty" validate:"omitempty,max=255"`
	Gender    *string `json:"gender,omitempty" validate:"omitempty,oneof=male female"`
	Phone     *string `json:"phone,omitempty" validate:"omitempty,max=50"`
}

// UserRepository mendefinisikan operasi database yang bisa dilakukan pada User
type UserRepository interface {
	FindAll() ([]User, error)
	FindByID(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	Save(user *User) (*User, error)
	Update(user *User) (*User, error)
	Delete(id uint) error
}

// UserService mendefinisikan logika bisnis yang bisa dilakukan pada User
type UserService interface {
	CreateUser(req UserCreateRequest) (*User, error)
	GetAllUsers() ([]User, error)
	GetUserByID(id uint) (*User, error)
	UpdateUser(id uint, req UserUpdateRequest) (*User, error)
	DeleteUser(id uint) error
}
