package user

import (
	"context"

	"github.com/aburizalpurnama/travel/internal/app/contract"
	"github.com/aburizalpurnama/travel/internal/app/model"
	"github.com/aburizalpurnama/travel/internal/app/payload"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo contract.UserRepository
}

// NewUserService membuat instance baru dari user service
func NewUserService(repo contract.UserRepository) contract.UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, req payload.CreateUserRequest) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	passwordHash := string(hashedPassword)
	user := &model.User{
		FirstName:    req.FirstName,
		FullName:     req.FullName,
		Gender:       req.Gender,
		Email:        &req.Email,
		Phone:        req.Phone,
		PasswordHash: &passwordHash,
		Role:         req.Role,
	}

	return s.repo.Save(ctx, user)
}

func (s *userService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	return s.repo.FindAll(ctx)
}

func (s *userService) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, id uint, req payload.UpdateUserRequest) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}

	if req.FullName != nil {
		user.FullName = *req.FullName
	}

	if req.Phone != nil {
		user.Phone = *req.Phone
	}

	return s.repo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
