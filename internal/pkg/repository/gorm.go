package repository

import (
	"context"
	"time"

	"github.com/aburizalpurnama/travel/internal/app/model"
	"github.com/aburizalpurnama/travel/internal/pkg/gormhelper"
	"github.com/aburizalpurnama/travel/internal/pkg/paginator"
	"gorm.io/gorm"
)

// M = Model, F = Filter
type GORM[M any, F any] struct {
	DB *gorm.DB
}

func NewGORM[T any, F any](db *gorm.DB) *GORM[T, F] {
	return &GORM[T, F]{DB: db}
}

func (r *GORM[M, F]) FindAll(ctx context.Context, page *int, size *int, filter *F) (data []M, err error) {
	query := r.DB.Where("deleted_on IS NULL")
	query, err = gormhelper.ParseFilter(query, filter)
	if err != nil {
		return nil, err
	}

	if page != nil && size != nil {
		offset := paginator.GetOffset(*page, *size)
		query.Offset(offset).Limit(*size)
	}

	err = query.Find(&data).Error
	return data, err
}

func (r *GORM[M, F]) Count(ctx context.Context, filter *model.ProductFilter) (count int64, err error) {
	var data M
	query := r.DB.Model(data).Where("deleted_on IS NULL")
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

func (r *GORM[M, F]) FindByID(ctx context.Context, id uint) (*M, error) {
	var data M
	err := r.DB.Where("deleted_on IS NULL").First(&data, id).Error
	return &data, err
}

func (r *GORM[M, F]) Save(ctx context.Context, data *M) (*M, error) {
	err := r.DB.Create(data).Error
	return data, err
}

func (r *GORM[M, F]) Update(ctx context.Context, data *M) (*M, error) {
	err := r.DB.Save(data).Error
	return data, err
}

func (r *GORM[M, F]) Delete(ctx context.Context, id uint) error {
	var data M
	return r.DB.Model(&data).Where("id = ?", id).Update("deleted_on", time.Now()).Error
}
