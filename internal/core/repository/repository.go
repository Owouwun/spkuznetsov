package repository

import (
	"context"

	"github.com/Owouwun/spkuznetsov/internal/core/logic/orders"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, req *orders.Order) (uint, error)
	UpdateOrder(ctx context.Context, id uint, req *orders.Order) error
	GetOrder(ctx context.Context, id uint) (*orders.Order, error)
}
