package payload

// CreateUserRequest adalah DTO untuk membuat user baru
type CreateUserRequest struct {
	FirstName string `json:"first_name" validate:"required,max=100"`
	FullName  string `json:"full_name" validate:"required,max=255"`
	Gender    string `json:"gender" validate:"required,oneof=male female"`
	Email     string `json:"email" validate:"required,email,max=320"`
	Phone     string `json:"phone" validate:"required,max=50"`
	Password  string `json:"password" validate:"required,min=8"`
	Role      string `json:"role" validate:"required,oneof=customer muthawif"`
}

// UpdateUserRequest adalah DTO untuk memperbarui user
type UpdateUserRequest struct {
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,max=100"`
	FullName  *string `json:"full_name,omitempty" validate:"omitempty,max=255"`
	Gender    *string `json:"gender,omitempty" validate:"omitempty,oneof=male female"`
	Phone     *string `json:"phone,omitempty" validate:"omitempty,max=50"`
}
