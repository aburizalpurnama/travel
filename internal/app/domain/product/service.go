package product

import (
	"context"
	"errors"

	"github.com/aburizalpurnama/travel/internal/app/contract"
	"github.com/aburizalpurnama/travel/internal/app/model"
	"github.com/aburizalpurnama/travel/internal/app/payload"
	"github.com/aburizalpurnama/travel/internal/pkg/apperror"
	"github.com/aburizalpurnama/travel/internal/pkg/dberror"
	"github.com/aburizalpurnama/travel/internal/pkg/response"
	"github.com/jinzhu/copier"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type service struct {
	uow contract.UnitOfWork
}

// NewService membuat instance baru dari product service
func NewService(uow contract.UnitOfWork) contract.ProductService {
	return &service{uow: uow}
}

func (s *service) CreateProduct(ctx context.Context, req payload.ProductCreateRequest) (*payload.ProductBaseResponse, error) {
	product := &model.Product{}
	err := copier.Copy(&product, &req)
	if err != nil {
		return nil, err
	}

	product.Price, err = decimal.NewFromString(req.Price)
	if err != nil {
		return nil, err
	}

	// Add required business logic here

	var created *model.Product
	err = s.uow.Execute(ctx, func(uow contract.UnitOfWork) error {
		var err error
		created, err = uow.Product().Save(ctx, product)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		if pgErr := dberror.GetError(err); pgErr != nil {
			switch pgErr.Code {
			case dberror.UniqueViolation:
				msg, details := dberror.ParseUniqueConstraintError(pgErr)
				return nil, apperror.New(apperror.DuplicateEntry, msg, err, details)
			}
		}

		return nil, apperror.New(apperror.Internal, "failed to save product", err, nil)
	}

	res := &payload.ProductBaseResponse{}
	err = copier.Copy(&res, created)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) GetAllProducts(ctx context.Context, req payload.ProductGetAllRequest) ([]payload.ProductBaseResponse, *response.Pagination, error) {
	var count int64
	var products []model.Product

	group, groupCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		var err error
		count, err = s.uow.Product().Count(groupCtx, req.ProductFilter)
		if err != nil {
			return err
		}

		return nil
	})

	group.Go(func() error {
		var err error
		products, err = s.uow.Product().FindAll(ctx, req.Page, req.Size, req.ProductFilter)
		if err != nil {
			return err
		}

		return nil
	})

	err := group.Wait()
	if err != nil {
		return nil, nil, err
	}

	var res []payload.ProductBaseResponse
	err = copier.Copy(&res, &products)
	if err != nil {
		return nil, nil, err
	}

	return res, response.NewPagination(req.Page, req.Size, &count), nil
}

func (s *service) GetProductByID(ctx context.Context, id uint) (*payload.ProductBaseResponse, error) {
	product, err := s.uow.Product().FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(apperror.ProductNotFound, "product not found", err, nil)
		}

		return nil, err
	}

	res := &payload.ProductBaseResponse{}
	err = copier.Copy(res, product)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) UpdateProduct(ctx context.Context, id uint, req payload.ProductUpdateRequest) (*payload.ProductBaseResponse, error) {
	product, err := s.uow.Product().FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = copier.CopyWithOption(&product, &req, copier.Option{IgnoreEmpty: true})
	if err != nil {
		return nil, err
	}

	// Add required business logic here

	updated, err := s.uow.Product().Update(ctx, product)
	if err != nil {
		return nil, err
	}

	res := &payload.ProductBaseResponse{}
	err = copier.Copy(res, updated)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) DeleteProduct(ctx context.Context, id uint) error {
	_, err := s.uow.Product().FindByID(ctx, id)
	if err != nil {
		return err
	}

	return s.uow.Product().Delete(ctx, id)
}
