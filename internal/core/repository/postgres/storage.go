package repository_postgres

import (
	"context"
	"database/sql"

	"github.com/seagumineko/spkuznetsov/internal/errors"
	"github.com/seagumineko/spkuznetsov/internal/core/logic/requests"
	"github.com/seagumineko/spkuznetsov/internal/core/repository"
)

type PostgresRequestRepository struct {
	db *sql.DB
}

// TODO: Add error handling
func NewRequestRepository(db *sql.DB) repository.RequestRepository {
	return &PostgresRequestRepository{db: db}
}

func (r *PostgresRequestRepository) CreateRequest(ctx context.Context, req *requests.Request) (int64, error) {
	return -1, core_errors.ErrNotImplemented
}

func (r *PostgresRequestRepository) UpdateRequest(ctx context.Context, id int64, req *requests.Request) error {
	return core_errors.ErrNotImplemented
}

func (r *PostgresRequestRepository) GetRequest(ctx context.Context, id int64) (*requests.Request, error) {
	return nil, core_errors.ErrNotImplemented
}
