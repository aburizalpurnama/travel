package user

import (
	"context"
	"time"

	"github.com/aburizalpurnama/travel/internal/app/contract"
	"github.com/aburizalpurnama/travel/internal/app/model"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository membuat instance baru dari GORM repository
func NewUserRepository(db *gorm.DB) contract.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll(ctx context.Context) ([]model.User, error) {
	var users []model.User
	err := r.db.Where("deleted_on IS NULL").Find(&users).Error
	return users, err
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.Where("deleted_on IS NULL").First(&user, id).Error
	return &user, err
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ? AND deleted_on IS NULL", email).First(&user).Error
	return &user, err
}

func (r *userRepository) Save(ctx context.Context, user *model.User) (*model.User, error) {
	err := r.db.Create(user).Error
	return user, err
}

func (r *userRepository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	err := r.db.Save(user).Error
	return user, err
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	// Soft delete
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("deleted_on", time.Now()).Error
}
