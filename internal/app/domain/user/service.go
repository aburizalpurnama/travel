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

func (s *userService) CreateUser(ctx context.Context, req payload.UserCreateRequest) (*payload.UserResponse, error) {
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

	created, err := s.repo.Save(ctx, user)
	if err != nil {
		return nil, err
	}

	res := &payload.UserResponse{}
	err = copier.Copy(res, created)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *userService) GetAllUsers(ctx context.Context, req payload.UserGetAllRequest) ([]payload.UserResponse, *response.Pagination, error) {
	var count int64
	var users []model.User
	var err error

	wg := sync.WaitGroup{}
	wg.Go(func() {
		count, err = s.repo.Count(ctx, req.UserFilter)
	})

	wg.Go(func() {
		users, err = s.repo.FindAll(ctx, req.Page, req.Size, req.UserFilter)
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

	if req.Page != nil && req.Size != nil {
		pagination.CurrentPage = req.Page
		pagination.PageSize = req.Size
		pagination.TotalItems = &count

		totalPages := int(count / int64(*req.Size))
		pagination.TotalPages = &totalPages
	}

	return resp, &pagination, nil
}

func (s *userService) GetUserByID(ctx context.Context, id uint) (*payload.UserResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	res := &payload.UserResponse{}
	err = copier.Copy(res, user)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *userService) UpdateUser(ctx context.Context, id uint, req payload.UserUpdateRequest) (*payload.UserResponse, error) {
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

	updated, err := s.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	res := &payload.UserResponse{}
	err = copier.Copy(res, updated)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
