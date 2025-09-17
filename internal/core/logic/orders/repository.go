package orders

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OrderRepository interface {
	GetOrders(ctx context.Context) ([]*Order, error)
	GetOrderByID(ctx context.Context, id uuid.UUID) (*Order, error)
	CreateOrder(ctx context.Context, order *Order) (uuid.UUID, error)
	UpdateOrder(ctx context.Context, order *Order) error
	Preschedule(ctx context.Context, id uuid.UUID, scheduledFor *time.Time) error
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

func (s *OrderService) GetOrders(ctx context.Context) ([]*Order, error) {
	return s.repo.GetOrders(ctx)
}

func (s *OrderService) Preschedule(ctx context.Context, id uuid.UUID, ScheduledFor *time.Time) error {
	return s.repo.Preschedule(ctx, id, ScheduledFor)
}
