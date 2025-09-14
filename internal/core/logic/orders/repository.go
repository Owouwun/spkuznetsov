package orders

import (
	"context"

	"github.com/google/uuid"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) (uuid.UUID, error)
	GetOrderByID(ctx context.Context, id uuid.UUID) (*Order, error)
	UpdateOrder(ctx context.Context, order *Order) error
}

type OrderService struct {
	repo OrderRepository
}

func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, pord *PrimaryOrder) (uuid.UUID, error) {
	order, err := pord.CreateNewOrder()
	if err != nil {
		return uuid.Nil, err
	}

	id, err := s.repo.CreateOrder(ctx, order)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id uuid.UUID) (*Order, error) {
	return s.repo.GetOrderByID(ctx, id)
}
