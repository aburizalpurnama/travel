package contract

import (
	"context"

	"github.com/aburizalpurnama/travel/internal/app/model"
	"github.com/aburizalpurnama/travel/internal/app/payload"
)

// UserService mendefinisikan logika bisnis yang bisa dilakukan pada model.User
type UserService interface {
	CreateUser(ctx context.Context, req payload.CreateUserRequest) (*model.User, error)
	GetAllUsers(ctx context.Context) ([]model.User, error)
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	UpdateUser(ctx context.Context, id uint, req payload.UpdateUserRequest) (*model.User, error)
	DeleteUser(ctx context.Context, id uint) error
}
