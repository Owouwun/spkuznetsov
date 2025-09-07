package repository_postgres

import (
	"context"
	"database/sql"

	"github.com/Owouwun/ipkuznetsov/internal/core"
	"github.com/Owouwun/ipkuznetsov/internal/core/logic/requests"
	"github.com/Owouwun/ipkuznetsov/internal/core/repository"
)

type PostgresRequestRepository struct {
	db *sql.DB
}

// TODO: Add error handling
func NewRequestRepository(db *sql.DB) repository.RequestRepository {
	return &PostgresRequestRepository{db: db}
}

func (r *PostgresRequestRepository) CreateRequest(ctx context.Context, req *requests.Request) (int64, error) {
	return -1, core.ErrNotImplemented
}

func (r *PostgresRequestRepository) UpdateRequest(ctx context.Context, id int64, req *requests.Request) error {
	return core.ErrNotImplemented
}

func (r *PostgresRequestRepository) GetRequest(ctx context.Context, id int64) (*requests.Request, error) {
	return nil, core.ErrNotImplemented
}
