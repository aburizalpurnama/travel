package user

import (
	"time"

	"github.com/aburizalpurnama/travel/internal/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository membuat instance baru dari GORM repository
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll() ([]domain.User, error) {
	var users []domain.User
	err := r.db.Where("deleted_on IS NULL").Find(&users).Error
	return users, err
}

func (r *userRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("deleted_on IS NULL").First(&user, id).Error
	return &user, err
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ? AND deleted_on IS NULL", email).First(&user).Error
	return &user, err
}

func (r *userRepository) Save(user *domain.User) (*domain.User, error) {
	err := r.db.Create(user).Error
	return user, err
}

func (r *userRepository) Update(user *domain.User) (*domain.User, error) {
	err := r.db.Save(user).Error
	return user, err
}

func (r *userRepository) Delete(id uint) error {
	// Soft delete
	return r.db.Model(&domain.User{}).Where("id = ?", id).Update("deleted_on", time.Now()).Error
}
