package repository_orders

import (
	"context"
	"time"

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

func (r *GormOrderRepository) getEntityByID(ctx context.Context, id uuid.UUID) (*entities.OrderEntity, error) {
	var orderEntity *entities.OrderEntity
	result := r.db.WithContext(ctx).
		Preload("Employee").
		First(&orderEntity, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return orderEntity, nil
}

func (r *GormOrderRepository) getEmployeeEntityByID(ctx context.Context, id uint) (*entities.EmployeeEntity, error) {
	var empEntity *entities.EmployeeEntity
	result := r.db.WithContext(ctx).
		First(&empEntity, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return empEntity, nil
}

func (r *GormOrderRepository) Create(ctx context.Context, ord *orders.Order) (uuid.UUID, error) {
	orderEntity := entities.NewOrderEntityFromLogic(ord)

	result := r.db.WithContext(ctx).Create(&orderEntity)
	if result.Error != nil {
		return uuid.Nil, result.Error
	}

	return orderEntity.ID, nil
}

func (r *GormOrderRepository) Update(ctx context.Context, ord *orders.Order) error {
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

func (r *GormOrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*orders.Order, error) {
	orderEntity, err := r.getEntityByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return orderEntity.ToLogicOrder(), nil
}

func (r *GormOrderRepository) GetAll(ctx context.Context) ([]*orders.Order, error) {
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

func (r *GormOrderRepository) Preschedule(ctx context.Context, id uuid.UUID, scheduledFor *time.Time) error {
	orderEntity, err := r.getEntityByID(ctx, id)
	if err != nil {
		return err
	}

	order := orderEntity.ToLogicOrder()

	err = order.Preschedule(scheduledFor)
	if err != nil {
		return err
	}

	orderEntity = entities.NewOrderEntityFromLogic(order)

	result := r.db.WithContext(ctx).
		Model(&orderEntity).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"Status":       orderEntity.Status,
			"ScheduledFor": scheduledFor,
		})

	return result.Error
}

func (r *GormOrderRepository) Assign(ctx context.Context, id uuid.UUID, empID uint) error {
	orderEntity, err := r.getEntityByID(ctx, id)
	if err != nil {
		return err
	}

	empEntity, err := r.getEmployeeEntityByID(ctx, empID)
	if err != nil {
		return err
	}

	order := orderEntity.ToLogicOrder()
	emp := empEntity.ToLogicEmployee()

	err = order.Assign(emp)
	if err != nil {
		return err
	}

	orderEntity = entities.NewOrderEntityFromLogic(order)

	result := r.db.WithContext(ctx).
		Model(&orderEntity).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"Status":     orderEntity.Status,
			"EmployeeID": orderEntity.Employee.ID,
		})

	return result.Error
}
