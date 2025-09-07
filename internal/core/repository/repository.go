package repository

import (
	"context"

	"github.com/Owouwun/ipkuznetsov/internal/core/logic/requests"
)

type RequestRepository interface {
	CreateRequest(ctx context.Context, req *requests.Request) (int64, error)
	UpdateRequest(ctx context.Context, id int64, req *requests.Request) error
	GetRequest(ctx context.Context, id int64) (*requests.Request, error)
}
