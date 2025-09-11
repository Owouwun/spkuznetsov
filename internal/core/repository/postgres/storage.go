package repository_postgres

import (
	"context"

	"github.com/Owouwun/spkuznetsov/internal/core/logic/orders"
	"github.com/Owouwun/spkuznetsov/internal/core/repository"
	deterrs "github.com/Owouwun/spkuznetsov/internal/errors"
	"gorm.io/gorm"
)

type GormOrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) repository.OrderRepository {
	return &GormOrderRepository{db: db}
}

func (r *GormOrderRepository) CreateOrder(ctx context.Context, req *orders.Order) (uint, error) {
	result := r.db.WithContext(ctx).Create(&req)
	if result.Error != nil {
		return 0, deterrs.NewDetErr(
			deterrs.QueryInsertFailed,
			deterrs.WithOriginalError(result.Error),
		)
	}
	return req.ID, nil
}

func (r *GormOrderRepository) UpdateOrder(ctx context.Context, id uint, req *orders.Order) error {
	req.ID = id
	result := r.db.WithContext(ctx).Save(&req)
	if result.Error != nil {
		return deterrs.NewDetErr(
			deterrs.QueryUpdateFailed,
			deterrs.WithOriginalError(result.Error),
		)
	}

	if result.RowsAffected == 0 {
		return deterrs.NewDetErr(deterrs.QuerySelectFailed)
	}

	return nil
}

func (r *GormOrderRepository) GetOrder(ctx context.Context, id uint) (*orders.Order, error) {
	var req orders.Order
	result := r.db.WithContext(ctx).Preload("Employee").First(&req, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, deterrs.NewDetErr(deterrs.QuerySelectFailed)
		}
		return nil, deterrs.NewDetErr(
			deterrs.QuerySelectFailed,
			deterrs.WithOriginalError(result.Error),
		)
	}
	return &req, nil
}
