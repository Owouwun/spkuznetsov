package repository

import (
	"context"

	"github.com/seagumineko/spkuznetsov/internal/core/logic/requests"
)

type RequestRepository interface {
	CreateRequest(ctx context.Context, req *requests.Request) (uint, error)
	UpdateRequest(ctx context.Context, id uint, req *requests.Request) error
	GetRequest(ctx context.Context, id uint) (*requests.Request, error)
}
