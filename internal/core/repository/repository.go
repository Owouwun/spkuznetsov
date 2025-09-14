package repository

import (
	"context"

	"github.com/Owouwun/spkuznetsov/internal/core/logic/orders"
	"github.com/google/uuid"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, req *orders.Order) (uuid.UUID, error)
	UpdateOrder(ctx context.Context, req *orders.Order) error
	GetOrder(ctx context.Context, id uuid.UUID) (*orders.Order, error)
}
