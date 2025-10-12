package contract

import (
	"context"

	"github.com/aburizalpurnama/travel/internal/app/model"
	"github.com/aburizalpurnama/travel/internal/app/payload"
	"github.com/aburizalpurnama/travel/internal/pkg/response"
)

// UserService mendefinisikan logika bisnis yang bisa dilakukan pada model.User
type UserService interface {
	CreateUser(ctx context.Context, req payload.UserCreateRequest) (*model.User, error)
	GetAllUsers(ctx context.Context, req payload.UserGetAllRequest) ([]payload.UserResponse, *response.Pagination, error)
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	UpdateUser(ctx context.Context, id uint, req payload.UserUpdateRequest) (*model.User, error)
	DeleteUser(ctx context.Context, id uint) error
}
