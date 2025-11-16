package product

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aburizalpurnama/travel/internal/app/contract"
	"github.com/aburizalpurnama/travel/internal/app/model"
	"github.com/aburizalpurnama/travel/internal/app/payload"
	"github.com/aburizalpurnama/travel/internal/pkg/apperror"
	"github.com/aburizalpurnama/travel/internal/pkg/dberror"
	"github.com/aburizalpurnama/travel/internal/pkg/response"
	"github.com/shopspring/decimal"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

var serviceTracer trace.Tracer = otel.Tracer("product.service")

type service struct {
	uow    contract.UnitOfWork
	mapper contract.Mapper
}

// NewService initializes a new instance of product service.
func NewService(uow contract.UnitOfWork, mapper contract.Mapper) contract.ProductService {
	return &service{uow: uow, mapper: mapper}
}

// CreateProduct handles the creation of a new product record.
func (s *service) CreateProduct(ctx context.Context, req payload.ProductCreateRequest) (*payload.ProductBaseResponse, error) {
	ctx, span := serviceTracer.Start(ctx, "CreateProduct")
	defer span.End()

	product := &model.Product{}
	err := s.mapper.ToModel(&req, product)
	if err != nil {
		return nil, apperror.New(apperror.Internal, "failed to map request", err, nil)
	}

	product.Price, err = decimal.NewFromString(req.Price)
	if err != nil {
		return nil, err
	}

	// TODO: Retrieve actor from context in production
	actor := struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}{ID: 1, Name: "John Doe"}

	actorJSON, _ := json.Marshal(actor)
	product.CreatedBy = actorJSON

	var created *model.Product
	err = s.uow.Execute(ctx, func(uowCtx context.Context, uow contract.UnitOfWork) error {
		created, err = uow.Product().Save(uowCtx, product)
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

		return nil, apperror.New(apperror.Internal, "failed to create product", err, nil)
	}

	res := &payload.ProductBaseResponse{}
	if err := s.mapper.ToResponse(created, res); err != nil {
		return nil, apperror.New(apperror.Internal, "failed to map response", err, nil)
	}

	return res, nil
}

// GetAllProducts retrieves a list of products with support for pagination and filtering.
func (s *service) GetAllProducts(ctx context.Context, req payload.ProductGetAllRequest) ([]payload.ProductBaseResponse, *response.Pagination, error) {
	ctx, span := serviceTracer.Start(ctx, "GetAllProducts")
	defer span.End()

	var count int64
	var products []model.Product

	// Use errgroup for concurrent data fetching (count and data)
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
		products, err = s.uow.Product().FindAll(groupCtx, req.Page, req.Size, req.ProductFilter)
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
	err = s.mapper.ToResponse(&products, &res)
	if err != nil {
		return nil, nil, err
	}

	return res, response.NewPagination(req.Page, req.Size, &count), nil
}

// GetProductByID retrieves a specific product by its unique identifier.
func (s *service) GetProductByID(ctx context.Context, id uint) (*payload.ProductBaseResponse, error) {
	ctx, span := serviceTracer.Start(ctx, "GetProductByID")
	defer span.End()

	product, err := s.uow.Product().FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(apperror.ProductNotFound, "product not found", err, nil)
		}
		return nil, err
	}

	res := &payload.ProductBaseResponse{}
	err = s.mapper.ToResponse(&product, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateProduct modifies an existing product's information.
func (s *service) UpdateProduct(ctx context.Context, id uint, req payload.ProductUpdateRequest) (*payload.ProductBaseResponse, error) {
	ctx, span := serviceTracer.Start(ctx, "UpdateProduct")
	defer span.End()

	product, err := s.uow.Product().FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = s.mapper.ToModel(&req, &product)
	if err != nil {
		return nil, err
	}

	updated, err := s.uow.Product().Update(ctx, product)
	if err != nil {
		return nil, err
	}

	res := &payload.ProductBaseResponse{}
	err = s.mapper.ToResponse(&updated, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteProduct removes a product record from the database.
func (s *service) DeleteProduct(ctx context.Context, id uint) error {
	ctx, span := serviceTracer.Start(ctx, "DeleteProduct")
	defer span.End()

	_, err := s.uow.Product().FindByID(ctx, id)
	if err != nil {
		return err
	}

	return s.uow.Product().Delete(ctx, id)
}
