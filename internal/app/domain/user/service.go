package user

import (
	"context"
	"sync"

	"github.com/aburizalpurnama/travel/internal/app/contract"
	"github.com/aburizalpurnama/travel/internal/app/model"
	"github.com/aburizalpurnama/travel/internal/app/payload"
	"github.com/aburizalpurnama/travel/internal/pkg/response"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo contract.UserRepository
}

// NewUserService membuat instance baru dari user service
func NewUserService(repo contract.UserRepository) contract.UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, req payload.UserCreateRequest) (*model.User, error) {
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

func (s *userService) GetAllUsers(ctx context.Context, req payload.UserGetAllRequest) ([]payload.UserResponse, *response.Pagination, error) {
	var err error
	var count int64
	var users []model.User

	wg := sync.WaitGroup{}
	wg.Go(func() {
		count, err = s.repo.Count(ctx, req.Filter)
	})

	wg.Go(func() {
		users, err = s.repo.FindAll(ctx, req.Option, req.Filter)
	})

	wg.Wait()
	if err != nil {
		return nil, nil, err
	}

	var resp []payload.UserResponse
	var pagination response.Pagination

	err = copier.Copy(&resp, &users)
	if err != nil {
		return nil, nil, err
	}

	if req.Option.Page != nil && req.Option.Size != nil {
		pagination.CurrentPage = req.Option.Page
		pagination.PageSize = req.Option.Size
		pagination.TotalItems = &count

		totalPages := int(count / int64(*req.Option.Size))
		pagination.TotalPages = &totalPages
	}

	return resp, &pagination, nil
}

func (s *userService) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, id uint, req payload.UserUpdateRequest) (*model.User, error) {
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
