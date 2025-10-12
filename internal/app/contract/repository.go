package contract

import (
	"context"

	"github.com/aburizalpurnama/travel/internal/app/model"
)

// UserRepository mendefinisikan operasi database yang bisa dilakukan pada model.User
type UserRepository interface {
	FindAll(ctx context.Context, page *int, size *int, filter *model.UserFilter) ([]model.User, error)
	Count(ctx context.Context, filter *model.UserFilter) (int64, error)
	FindByID(ctx context.Context, id uint) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Save(ctx context.Context, user *model.User) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, id uint) error
}
