package orders

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OrderRepository interface {
	GetAll(ctx context.Context) ([]*Order, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Order, error)
	Create(ctx context.Context, order *Order) (uuid.UUID, error)
	Update(ctx context.Context, order *Order) error
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

func (s *OrderService) Create(ctx context.Context, pord *PrimaryOrder) (uuid.UUID, error) {
	order, err := pord.CreateNewOrder()
	if err != nil {
		return uuid.Nil, err
	}

	id, err := s.repo.Create(ctx, order)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (s *OrderService) GetByID(ctx context.Context, id uuid.UUID) (*Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *OrderService) GetAll(ctx context.Context) ([]*Order, error) {
	return s.repo.GetAll(ctx)
}

func (s *OrderService) Preschedule(ctx context.Context, id uuid.UUID, ScheduledFor *time.Time) error {
	return s.repo.Preschedule(ctx, id, ScheduledFor)
}
