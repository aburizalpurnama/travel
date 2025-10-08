package user

import (
	"github.com/aburizalpurnama/travel/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo domain.UserRepository
}

// NewUserService membuat instance baru dari user service
func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(req domain.UserCreateRequest) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	passwordHash := string(hashedPassword)
	user := &domain.User{
		FirstName:    req.FirstName,
		FullName:     req.FullName,
		Gender:       req.Gender,
		Email:        &req.Email,
		Phone:        req.Phone,
		PasswordHash: &passwordHash,
		Role:         req.Role,
	}

	return s.repo.Save(user)
}

func (s *userService) GetAllUsers() ([]domain.User, error) {
	return s.repo.FindAll()
}

func (s *userService) GetUserByID(id uint) (*domain.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) UpdateUser(id uint, req domain.UserUpdateRequest) (*domain.User, error) {
	user, err := s.repo.FindByID(id)
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

	return s.repo.Update(user)
}

func (s *userService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}
