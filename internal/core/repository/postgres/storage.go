package repository_postgres

import (
	"context"

	"github.com/seagumineko/spkuznetsov/internal/core/logic/requests"
	"github.com/seagumineko/spkuznetsov/internal/core/repository"
	deterrs "github.com/seagumineko/spkuznetsov/internal/errors"
	"gorm.io/gorm"
)

type GormRequestRepository struct {
	db *gorm.DB
}

func NewRequestRepository(db *gorm.DB) repository.RequestRepository {
	return &GormRequestRepository{db: db}
}

func (r *GormRequestRepository) CreateRequest(ctx context.Context, req *requests.Request) (uint, error) {
	result := r.db.WithContext(ctx).Create(&req)
	if result.Error != nil {
		return 0, deterrs.NewDetErr(
			deterrs.QueryInsertFailed,
			deterrs.WithOriginalError(result.Error),
		)
	}
	return req.ID, nil
}

func (r *GormRequestRepository) UpdateRequest(ctx context.Context, id uint, req *requests.Request) error {
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

func (r *GormRequestRepository) GetRequest(ctx context.Context, id uint) (*requests.Request, error) {
	var req requests.Request
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
