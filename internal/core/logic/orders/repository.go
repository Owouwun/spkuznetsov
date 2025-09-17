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
	Assign(ctx context.Context, id uuid.UUID, empID uint) error
	Schedule(ctx context.Context, id uuid.UUID, scheduledFor *time.Time) error
	Progress(ctx context.Context, id uuid.UUID, empDescr string) error
	Complete(ctx context.Context, id uuid.UUID) error
	Close(ctx context.Context, id uuid.UUID) error
	Cancel(ctx context.Context, id uuid.UUID, reason string) error
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

func (s *OrderService) Preschedule(ctx context.Context, id uuid.UUID, scheduledFor *time.Time) error {
	return s.repo.Preschedule(ctx, id, scheduledFor)
}

func (s *OrderService) Assign(ctx context.Context, id uuid.UUID, empID uint) error {
	return s.repo.Assign(ctx, id, empID)
}

func (s *OrderService) Schedule(ctx context.Context, id uuid.UUID, scheduledFor *time.Time) error {
	return s.repo.Schedule(ctx, id, scheduledFor)
}

func (s *OrderService) Progress(ctx context.Context, id uuid.UUID, empDescr string) error {
	return s.repo.Progress(ctx, id, empDescr)
}

func (s *OrderService) Complete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Complete(ctx, id)
}

func (s *OrderService) Close(ctx context.Context, id uuid.UUID) error {
	return s.repo.Close(ctx, id)
}

func (s *OrderService) Cancel(ctx context.Context, id uuid.UUID, reason string) error {
	return s.repo.Cancel(ctx, id, reason)
}
