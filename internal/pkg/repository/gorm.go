package repository

import (
	"context"
	"time"

	"github.com/aburizalpurnama/travel/internal/app/model"
	"github.com/aburizalpurnama/travel/internal/pkg/gormhelper"
	"github.com/aburizalpurnama/travel/internal/pkg/paginator"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

var gormTracer trace.Tracer = otel.Tracer("repository.gorm")

// GORM is a generic repository implementation using GORM.
// M represents the Model type, and F represents the Filter type.
type GORM[M any, F any] struct {
	db *gorm.DB
}

// NewGORM creates a new instance of the generic GORM repository.
func NewGORM[M any, F any](db *gorm.DB) *GORM[M, F] {
	return &GORM[M, F]{db: db}
}

// FindAll retrieves a list of records based on pagination and filter criteria.
func (r *GORM[M, F]) FindAll(ctx context.Context, page *int, size *int, filter *F) (data []M, err error) {
	ctx, span := gormTracer.Start(ctx, "GORM.FindAll")
	defer span.End()

	query := r.db.WithContext(ctx).Where("deleted_on IS NULL")

	// Apply dynamic filtering based on the filter struct
	query, err = gormhelper.ParseFilter(query, filter)
	if err != nil {
		return nil, err
	}

	// Apply pagination if page and size are provided
	if page != nil && size != nil {
		offset := paginator.GetOffset(*page, *size)
		query = query.Offset(offset).Limit(*size)
	}

	err = query.Find(&data).Error
	return data, err
}

// Count retrieves the total number of records matching the filter criteria.
func (r *GORM[M, F]) Count(ctx context.Context, filter *model.ProductFilter) (count int64, err error) {
	ctx, span := gormTracer.Start(ctx, "GORM.Count")
	defer span.End()

	var data M
	query := r.db.WithContext(ctx).Model(data).Where("deleted_on IS NULL")

	// Apply dynamic filtering
	query, err = gormhelper.ParseFilter(query, filter)
	if err != nil {
		return 0, err
	}

	err = query.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

// FindByID retrieves a single record by its unique identifier (ID).
func (r *GORM[M, F]) FindByID(ctx context.Context, id uint) (*M, error) {
	ctx, span := gormTracer.Start(ctx, "GORM.FindByID")
	defer span.End()

	var data M
	err := r.db.WithContext(ctx).Where("deleted_on IS NULL").First(&data, id).Error
	return &data, err
}

// Save persists a new record to the database.
func (r *GORM[M, F]) Save(ctx context.Context, data *M) (*M, error) {
	ctx, span := gormTracer.Start(ctx, "GORM.Save")
	defer span.End()

	err := r.db.WithContext(ctx).Create(data).Error
	return data, err
}

// Update modifies an existing record in the database.
func (r *GORM[M, F]) Update(ctx context.Context, data *M) (*M, error) {
	ctx, span := gormTracer.Start(ctx, "GORM.Update")
	defer span.End()

	err := r.db.WithContext(ctx).Save(data).Error
	return data, err
}

// Delete performs a soft delete on a record by setting the deleted_on timestamp.
func (r *GORM[M, F]) Delete(ctx context.Context, id uint) error {
	ctx, span := gormTracer.Start(ctx, "GORM.Delete")
	defer span.End()

	var data M
	// Perform a soft delete by updating the deleted_on column
	return r.db.WithContext(ctx).Model(&data).Where("id = ?", id).Update("deleted_on", time.Now()).Error
}
