package product

import (
	"context"
	"time"

	"github.com/aburizalpurnama/travel/internal/app/contract"
	"github.com/aburizalpurnama/travel/internal/app/model"
	"github.com/aburizalpurnama/travel/internal/pkg/gormhelper"
	"github.com/aburizalpurnama/travel/internal/pkg/paginator"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository membuat instance baru dari GORM repository
func NewProductRepository(db *gorm.DB) contract.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindAll(ctx context.Context, page *int, size *int, filter *model.ProductFilter) (products []model.Product, err error) {
	// .Where("deleted_on IS NULL")
	query := r.db
	query, err = gormhelper.ParseFilter(query, filter)
	if err != nil {
		return nil, err
	}

	if page != nil && size != nil {
		offset := paginator.GetOffset(*page, *size)
		query.Offset(offset).Limit(*size)
	}

	err = query.Find(&products).Error
	return products, err
}

func (r *productRepository) Count(ctx context.Context, filter *model.ProductFilter) (count int64, err error) {
	// .Where("deleted_on IS NULL")
	query := r.db.Model(&model.Product{})
	query, err = gormhelper.ParseFilter(query, filter)
	if err != nil {
		return 0, err
	}

	err = query.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return
}

func (r *productRepository) FindByID(ctx context.Context, id uint) (*model.Product, error) {
	var product model.Product
	err := r.db.Where("deleted_on IS NULL").First(&product, id).Error
	return &product, err
}

func (r *productRepository) Save(ctx context.Context, product *model.Product) (*model.Product, error) {
	err := r.db.Create(product).Error
	return product, err
}

func (r *productRepository) Update(ctx context.Context, product *model.Product) (*model.Product, error) {
	err := r.db.Save(product).Error
	return product, err
}

func (r *productRepository) Delete(ctx context.Context, id uint) error {
	return r.db.Model(&model.Product{}).Where("id = ?", id).Update("deleted_on", time.Now()).Error
}
