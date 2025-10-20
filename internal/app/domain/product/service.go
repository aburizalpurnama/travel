package product

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aburizalpurnama/travel/internal/app/contract"
	"github.com/aburizalpurnama/travel/internal/app/model"
	"github.com/aburizalpurnama/travel/internal/app/payload"
	"github.com/aburizalpurnama/travel/internal/pkg/apperror"
	"github.com/aburizalpurnama/travel/internal/pkg/response"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jinzhu/copier"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type productService struct {
	repo contract.ProductRepository
}

// NewProductService membuat instance baru dari product service
func NewProductService(repo contract.ProductRepository) contract.ProductService {
	return &productService{repo: repo}
}

func (s *productService) CreateProduct(ctx context.Context, req payload.ProductCreateRequest) (*payload.ProductBaseResponse, error) {
	product := &model.Product{}
	err := copier.Copy(&product, &req)
	if err != nil {
		return nil, err
	}

	actor := struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}{ID: 1, Name: "John Doe"}

	actorJSON, _ := json.Marshal(actor)
	product.CreatedBy = actorJSON

	product.Price, err = decimal.NewFromString(req.Price)
	if err != nil {
		return nil, err
	}

	// Add required business logic here

	created, err := s.repo.Save(ctx, product)
	if err != nil {
		// TODO: handle ketika ada error db seperti duplicate entry atau yang lain.
		// check apakah di table terdapat unique constraint selain primary key
		var pgErr *pgconn.PgError

		// Coba 'cast' error-nya ke PgError
		if errors.As(err, &pgErr) {
			fmt.Printf("pgErr: %+v\n", *pgErr)

			// Periksa KODE SQLSTATE untuk 'unique_violation'
			if pgErr.Code == "23505" {

				// TODO: define semua based on unique constraints
				// (Opsional tapi sangat baik) Cek nama constraint-nya
				switch pgErr.ConstraintName {
				case "products_name_unique":
					return nil, apperror.New(apperror.ProductNameExists, "product name already exists", err)
				case "ux_products_uid_active":
					return nil, apperror.New(apperror.DuplicateEntry, "product UID already exists", err)
				default:
					return nil, apperror.New(apperror.DuplicateEntry, "unique constraint violated", err)
				}
			}
		}

		return nil, err
	}

	res := &payload.ProductBaseResponse{}
	err = copier.Copy(&res, created)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *productService) GetAllProducts(ctx context.Context, req payload.ProductGetAllRequest) ([]payload.ProductBaseResponse, *response.Pagination, error) {
	var count int64
	var products []model.Product

	group, groupCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		var err error
		count, err = s.repo.Count(groupCtx, req.ProductFilter)
		if err != nil {
			return err
		}

		return nil
	})

	group.Go(func() error {
		var err error
		products, err = s.repo.FindAll(ctx, req.Page, req.Size, req.ProductFilter)
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

func (s *productService) GetProductByID(ctx context.Context, id uint) (*payload.ProductBaseResponse, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(apperror.ProductNotFound, "product not found", err)
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

func (s *productService) UpdateProduct(ctx context.Context, id uint, req payload.ProductUpdateRequest) (*payload.ProductBaseResponse, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = copier.CopyWithOption(&product, &req, copier.Option{IgnoreEmpty: true})
	if err != nil {
		return nil, err
	}

	// Add required business logic here

	updated, err := s.repo.Update(ctx, product)
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

func (s *productService) DeleteProduct(ctx context.Context, id uint) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}
