package user

import (
	"context"
	"time"

	"github.com/aburizalpurnama/travel/internal/app/contract"
	"github.com/aburizalpurnama/travel/internal/app/model"
	"github.com/aburizalpurnama/travel/internal/pkg/paginator"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository membuat instance baru dari GORM repository
func NewUserRepository(db *gorm.DB) contract.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll(ctx context.Context, opts paginator.OffsetBasedOption, filter model.UserFilter) (users []model.User, err error) {
	db := r.db.Where("deleted_on IS NULL")

	if filter.ID != nil {
		db.Where("id", filter.ID)
	}

	if filter.UID != nil {
		db.Where("uid", filter.UID)
	}

	if filter.Email != nil {
		db.Where("email", filter.Email)
	}

	if filter.IsActive != nil {
		db.Where("is_active", filter.IsActive)
	}

	if filter.IsVerified != nil {
		db.Where("is_verified", filter.IsVerified)
	}

	if opts.Page != nil && opts.Size != nil {
		offset := paginator.GetOffset(*opts.Page, *opts.Size)
		db.Offset(offset).Limit(*opts.Size)
	}

	err = db.Find(&users).Error
	return users, err
}

func (r *userRepository) Count(ctx context.Context, filter model.UserFilter) (count int64, err error) {
	db := r.db.Model(&model.User{}).Where("deleted_on IS NULL")

	if filter.ID != nil {
		db.Where("id", filter.ID)
	}

	if filter.UID != nil {
		db.Where("uid", filter.UID)
	}

	if filter.Email != nil {
		db.Where("email", filter.Email)
	}

	if filter.IsActive != nil {
		db.Where("is_active", filter.IsActive)
	}

	if filter.IsVerified != nil {
		db.Where("is_verified", filter.IsVerified)
	}

	err = db.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return
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
