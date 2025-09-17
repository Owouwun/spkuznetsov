package repository_orders

import (
	"context"

	"github.com/Owouwun/spkuznetsov/internal/core/logic/orders"
	"github.com/Owouwun/spkuznetsov/internal/core/repository/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormOrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) orders.OrderRepository {
	return &GormOrderRepository{db: db}
}

func (r *GormOrderRepository) CreateOrder(ctx context.Context, ord *orders.Order) (uuid.UUID, error) {
	orderEntity := entities.NewOrderEntityFromLogic(ord)

	result := r.db.WithContext(ctx).Create(&orderEntity)
	if result.Error != nil {
		return uuid.Nil, result.Error
	}

	return orderEntity.ID, nil
}

func (r *GormOrderRepository) UpdateOrder(ctx context.Context, ord *orders.Order) error {
	orderEntity := entities.NewOrderEntityFromLogic(ord)

	result := r.db.WithContext(ctx).
		Model(&orderEntity).
		Where("id = ?", ord.ID).
		Select(
			"ClientName",
			"ClientPhone",
			"Address",
			"ClientDescription",
			"EmployeeID",
			"CancelReason",
			"Status",
			"EmployeeDescription",
			"ScheduledFor",
		).
		Updates(orderEntity)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *GormOrderRepository) GetOrderByID(ctx context.Context, id uuid.UUID) (*orders.Order, error) {
	var orderEntity entities.OrderEntity
	result := r.db.WithContext(ctx).
		Preload("Employee").
		First(&orderEntity, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return orderEntity.ToLogicOrder(), nil
}

func (r *GormOrderRepository) GetOrders(ctx context.Context) ([]*orders.Order, error) {
	var orderEntities []entities.OrderEntity
	result := r.db.WithContext(ctx).
		Preload("Employee").
		Find(&orderEntities)

	if result.Error != nil {
		return nil, result.Error
	}

	var logicOrders []*orders.Order
	for _, entity := range orderEntities {
		logicOrders = append(logicOrders, entity.ToLogicOrder())
	}

	return logicOrders, nil
}
