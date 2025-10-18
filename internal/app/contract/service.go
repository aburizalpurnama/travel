package contract

import (
	"context"

	"github.com/aburizalpurnama/travel/internal/app/payload"
	"github.com/aburizalpurnama/travel/internal/pkg/response"
)

// UserService mendefinisikan logika bisnis yang bisa dilakukan pada model.User
type UserService interface {
	CreateUser(ctx context.Context, req payload.UserCreateRequest) (*payload.UserResponse, error)
	GetAllUsers(ctx context.Context, req payload.UserGetAllRequest) ([]payload.UserResponse, *response.Pagination, error)
	GetUserByID(ctx context.Context, id uint) (*payload.UserResponse, error)
	UpdateUser(ctx context.Context, id uint, req payload.UserUpdateRequest) (*payload.UserResponse, error)
	DeleteUser(ctx context.Context, id uint) error
}

// ProductService mendefinisikan logika bisnis yang bisa dilakukan pada model.Product
type ProductService interface {
	CreateProduct(ctx context.Context, req payload.ProductCreateRequest) (*payload.ProductBaseResponse, error)
	GetAllProducts(ctx context.Context, req payload.ProductGetAllRequest) ([]payload.ProductBaseResponse, *response.Pagination, error)
	GetProductByID(ctx context.Context, id uint) (*payload.ProductBaseResponse, error)
	UpdateProduct(ctx context.Context, id uint, req payload.ProductUpdateRequest) (*payload.ProductBaseResponse, error)
	DeleteProduct(ctx context.Context, id uint) error
}
